package api

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/openshift/origin/pkg/cmd/util/clientcmd"

	k8client "k8s.io/kubernetes/pkg/client"
	kclientcmd "k8s.io/kubernetes/pkg/client/clientcmd"
	kclientcmdapi "k8s.io/kubernetes/pkg/client/clientcmd/api"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/kubectl/resource"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/pmylund/go-cache"
	"github.com/spf13/pflag"
)

// ClientWrapper wraps an OpenShift client
// and supports caching
type ClientWrapper struct {
	factory         *clientcmd.Factory
	namespacesCache *cache.Cache
	resourcesCache  *cache.Cache
}

// NewClientWrapper build a new ClientWrapper instance
// with or without caching
func NewClientWrapper(withCache bool) *ClientWrapper {
	factory := getFactory()

	var resourcesCache *cache.Cache
	if withCache {
		resourcesCache = cache.New(5*time.Minute, 30*time.Second)
	} else {
		resourcesCache = nil
	}

	return &ClientWrapper{
		factory:         factory,
		resourcesCache:  resourcesCache,
		namespacesCache: cache.New(5*time.Minute, 30*time.Second),
	}
}

// LoadData build a Data instance, populated with data for the given resource types.
// You can use ResourceTypeAll to get data for all resources types.
// If caching is enabled, it will use the cache if there are fresh data in it.
func (cw *ClientWrapper) LoadData(resourceTypes ...ResourceType) (*Data, error) {
	data := &Data{}

	namespaces, err := cw.GetAvailableNamespaces()
	if err != nil {
		return nil, err
	}

	channels := []<-chan *DataWrapper{}
	for _, resourceType := range resourceTypes {
		switch resourceType {
		case ResourceTypeApplication:
			// nothing to load: we will extract the applications later...
		case ResourceTypeContainer:
			// nothing to load: we will extract the containers from the pods later...
		case ResourceTypeProject:
			channels = append(channels, cw.AsyncListResources(resourceType, "openshift"))
		default:
			channels = append(channels, cw.AsyncListResources(resourceType, namespaces...))
		}
	}

	for _, channel := range channels {
		select {
		case d := <-channel:
			if d.Errors != nil {
				return nil, fmt.Errorf("Failed to load data: %v", d.Errors)
			}
			data.Merge(d.Data)
		case <-time.After(10 * time.Second):
			return nil, fmt.Errorf("Timed out while loading data!")
		}
	}

	for _, resourceType := range resourceTypes {
		switch resourceType {
		case ResourceTypePod:
			data.RemoveBuilderAndDeployerPods()
		case ResourceTypeContainer:
			data.ExtractContainersFromPods()
		case ResourceTypeApplication:
			data.ExtractApplicationsFromDeploymentConfigs()
		}
	}

	return data, nil
}

// GetAvailableNamespaces retrieves all available namespaces.
func (cw *ClientWrapper) GetAvailableNamespaces() ([]string, error) {
	if namespaces, found := cw.namespacesCache.Get("namespaces"); found {
		return namespaces.([]string), nil
	}

	client, _, err := cw.factory.Clients()
	if err != nil {
		return nil, err
	}

	projectList, err := client.Projects().List(labels.Everything(), fields.Everything())
	if err != nil {
		return nil, err
	}

	namespaces := []string{}
	for _, project := range projectList.Items {
		namespaces = append(namespaces, project.Name)
	}

	cw.namespacesCache.Set("namespaces", namespaces, cache.DefaultExpiration)
	return namespaces, nil
}

// AsyncListResources retrieves the list of resources for the given resource type.
// You can restrict the resources to one or more namespaces (by default, it will use all available namespaces).
// It works asynchronously, and returns a channel immediately.
// When the result will be available, it will be send to the channel, which will then be closed.
// The result is a DataWrapper instance, that contains the data or an error.
func (cw *ClientWrapper) AsyncListResources(resourceType ResourceType, namespaces ...string) <-chan *DataWrapper {
	c := make(chan *DataWrapper)
	go func() {
		resources, errs := cw.ListResources(resourceType, namespaces...)
		data := &DataWrapper{
			Data: &Data{},
		}
		if errs != nil {
			data.Errors = errs
		} else {
			data.Set(resourceType, resources)
		}
		c <- data
		close(c)
	}()
	return c
}

// ListResources retrieves the list of resources for the given resource type.
// You can restrict the resources to one or more namespaces (by default, it will use all available namespaces).
// It returns the resources as a slice of interface{}, or the errors.
func (cw *ClientWrapper) ListResources(resourceType ResourceType, namespaces ...string) ([]interface{}, []error) {
	if cw.resourcesCache != nil {
		if resources, found := cw.resourcesCache.Get(string(resourceType)); found {
			return resources.([]interface{}), nil
		}
	}

	helper, version, err := cw.getHelperForResource(resourceType)
	if err != nil {
		return nil, []error{err}
	}

	if len(namespaces) == 0 {
		namespaces, err = cw.GetAvailableNamespaces()
		if err != nil {
			return nil, []error{err}
		}
	}

	results := []interface{}{}
	errs := []error{}
	for _, namespace := range namespaces {
		result, err := helper.List(namespace, version, labels.Everything())
		if err != nil {
			errs = append(errs, err)
			continue
		}

		items, err := extractItems(result)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		results = append(results, items...)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	if cw.resourcesCache != nil {
		cw.resourcesCache.Set(string(resourceType), results, cache.DefaultExpiration)
	}
	return results, nil
}

// getHelperForResource builds an API resource Helper configured to works with the given resource type.
// It returns the helper and the API version to use, or an error
func (cw *ClientWrapper) getHelperForResource(resourceType ResourceType) (helper *resource.Helper, version string, err error) {
	mapper, _ := cw.factory.Object()
	version, kind, err := mapper.VersionAndKindForResource(string(resourceType))
	if err != nil {
		return
	}

	mapping, err := mapper.RESTMapping(kind, version)
	if err != nil {
		return
	}

	client, err := cw.factory.RESTClient(mapping)
	if err != nil {
		return
	}

	helper = resource.NewHelper(client, mapping)
	return
}

// DataWrapper wraps a Data instance and a slice or errors that may have happened while loading the data
type DataWrapper struct {
	*Data
	Errors []error
}

// extractItems extracts the Items field value from the given k8s API Object.
// It is based on the assumption that the given API Object is a List of Items.
func extractItems(obj runtime.Object) ([]interface{}, error) {
	objValue := reflect.ValueOf(obj).Elem()
	if !objValue.IsValid() || objValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Failed to extract Items from object %T: it is not a Struct!", obj)
	}

	itemsField := objValue.FieldByName("Items")
	if !itemsField.IsValid() || !itemsField.CanInterface() {
		return nil, fmt.Errorf("Failed to extract Items from object %T: it has no 'Items' field!", obj)
	}

	items := itemsField.Interface()
	itemsValue := reflect.ValueOf(items)
	if !itemsValue.IsValid() || itemsValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Failed to extract Items from object %T: it is not a Slice!", items)
	}

	results := []interface{}{}
	for i := 0; i < itemsValue.Len(); i++ {
		valueValue := itemsValue.Index(i)
		if !valueValue.IsValid() || !valueValue.CanInterface() {
			return nil, fmt.Errorf("Failed to extract Item %v from object %T: it is not valid!", i, items)
		}

		value := valueValue.Interface()
		results = append(results, value)
	}
	return results, nil
}

// getFactory returns an OpenShift's Factory
// It first tries to use the config that is made available when we are running in a cluster
// and then fallback to a standard factory (using the default config files)
func getFactory() *clientcmd.Factory {
	factory, err := getFactoryFromCluster()
	if err != nil {
		log.Printf("Seems like we are not running in an OpenShift environment (%s), falling back to building a std factory...", err)
		factory = clientcmd.New(pflag.NewFlagSet("openshift-factory", pflag.ContinueOnError))
	}

	return factory
}

// getFactoryFromCluster returns an OpenShift's Factory
// using the config that is made available when we are running in a cluster
// (using environment variables and token secret file)
// or an error if those are not available (meaning we are not running in a cluster)
func getFactoryFromCluster() (*clientcmd.Factory, error) {
	clusterConfig, err := k8client.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// keep only what we need to initialize a factory
	overrides := &kclientcmd.ConfigOverrides{
		ClusterInfo: kclientcmdapi.Cluster{
			Server:     clusterConfig.Host,
			APIVersion: clusterConfig.Version,
		},
		AuthInfo: kclientcmdapi.AuthInfo{
			Token: clusterConfig.BearerToken,
		},
		Context: kclientcmdapi.Context{},
	}

	if len(clusterConfig.TLSClientConfig.CAFile) > 0 {
		// FIXME "x509: cannot validate certificate for x.x.x.x because it doesn't contain any IP SANs"
		// overrides.ClusterInfo.CertificateAuthority = clusterConfig.TLSClientConfig.CAFile
		overrides.ClusterInfo.InsecureSkipTLSVerify = true
	} else {
		overrides.ClusterInfo.InsecureSkipTLSVerify = true
	}

	config := kclientcmd.NewDefaultClientConfig(*kclientcmdapi.NewConfig(), overrides)

	factory := clientcmd.NewFactory(config)
	return factory, nil
}
