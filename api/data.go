package api

import (
	"fmt"
	"sort"

	buildapi "github.com/openshift/origin/pkg/build/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	imageapi "github.com/openshift/origin/pkg/image/api"
	projectapi "github.com/openshift/origin/pkg/project/api"
	routeapi "github.com/openshift/origin/pkg/route/api"

	kapi "k8s.io/kubernetes/pkg/api"
)

// Data contains all available data retrieved from the API
type Data struct {
	Applications           []Application
	Projects               []projectapi.Project
	Routes                 []routeapi.Route
	Services               []kapi.Service
	Pods                   []kapi.Pod
	Containers             []kapi.Container
	ImageStreams           []imageapi.ImageStream
	BuildConfigs           []buildapi.BuildConfig
	Builds                 []buildapi.Build
	DeploymentConfigs      []deployapi.DeploymentConfig
	ReplicationControllers []kapi.ReplicationController
	Events                 []kapi.Event
}

// Merge merges the given Data instances in this instance
func (d *Data) Merge(others ...*Data) {
	for _, other := range others {
		if d.Projects == nil && other.Projects != nil {
			d.Projects = other.Projects
		}
		if d.Routes == nil && other.Routes != nil {
			d.Routes = other.Routes
		}
		if d.Services == nil && other.Services != nil {
			d.Services = other.Services
		}
		if d.Pods == nil && other.Pods != nil {
			d.Pods = other.Pods
		}
		if d.Containers == nil && other.Containers != nil {
			d.Containers = other.Containers
		}
		if d.ImageStreams == nil && other.ImageStreams != nil {
			d.ImageStreams = other.ImageStreams
		}
		if d.BuildConfigs == nil && other.BuildConfigs != nil {
			d.BuildConfigs = other.BuildConfigs
		}
		if d.Builds == nil && other.Builds != nil {
			d.Builds = other.Builds
		}
		if d.DeploymentConfigs == nil && other.DeploymentConfigs != nil {
			d.DeploymentConfigs = other.DeploymentConfigs
		}
		if d.ReplicationControllers == nil && other.ReplicationControllers != nil {
			d.ReplicationControllers = other.ReplicationControllers
		}
		if d.Events == nil && other.Events != nil {
			d.Events = other.Events
		}
	}
}

// Set adds resources to this Data instance
func (d *Data) Set(resourceType ResourceType, resources []interface{}) error {
	switch resourceType {

	case ResourceTypeProject:
		return d.SetProjects(resources)

	case ResourceTypeRoute:
		return d.SetRoutes(resources)

	case ResourceTypeService:
		return d.SetServices(resources)

	case ResourceTypePod:
		return d.SetPods(resources)

	case ResourceTypeContainer:
		return d.SetContainers(resources)

	case ResourceTypeImageStream:
		return d.SetImageStreams(resources)

	case ResourceTypeBuildConfig:
		return d.SetBuildConfigs(resources)

	case ResourceTypeBuild:
		return d.SetBuilds(resources)

	case ResourceTypeDeploymentConfig:
		return d.SetDeploymentConfigs(resources)

	case ResourceTypeReplicationController:
		return d.SetReplicationControllers(resources)

	case ResourceTypeEvent:
		return d.SetEvents(resources)

	default:
		return fmt.Errorf("Unknown resource type %v!", resourceType)
	}
}

func (d *Data) SetProjects(projects []interface{}) error {
	d.Projects = []projectapi.Project{}
	for _, obj := range projects {
		project, ok := obj.(projectapi.Project)
		if !ok {
			return fmt.Errorf("Wrong type %T for projects!", projects)
		}
		d.Projects = append(d.Projects, project)
	}
	return nil
}

func (d *Data) SetRoutes(routes []interface{}) error {
	d.Routes = []routeapi.Route{}
	for _, obj := range routes {
		route, ok := obj.(routeapi.Route)
		if !ok {
			return fmt.Errorf("Wrong type %T for routes!", routes)
		}
		d.Routes = append(d.Routes, route)
	}
	return nil
}

func (d *Data) SetServices(services []interface{}) error {
	d.Services = []kapi.Service{}
	for _, obj := range services {
		svc, ok := obj.(kapi.Service)
		if !ok {
			return fmt.Errorf("Wrong type %T for services!", services)
		}
		d.Services = append(d.Services, svc)
	}
	return nil
}

func (d *Data) SetPods(pods []interface{}) error {
	d.Pods = []kapi.Pod{}
	for _, obj := range pods {
		pod, ok := obj.(kapi.Pod)
		if !ok {
			return fmt.Errorf("Wrong type %T for pods!", pods)
		}
		d.Pods = append(d.Pods, pod)
	}
	return nil
}

// RemoveBuilderAndDeployerPods removes the builders and deployers pods
// from the pods stored in this Data instance.
// It returns the builders and deployers pods that have been removed.
func (d *Data) RemoveBuilderAndDeployerPods() (builderPods []kapi.Pod, deployerPods []kapi.Pod) {
	pods := d.Pods
	d.Pods = pods[:0]
	for _, pod := range pods {
		if _, found := pod.Labels[buildapi.BuildLabel]; found {
			builderPods = append(builderPods, pod)
		} else if _, found := pod.Labels[deployapi.DeployerPodForDeploymentLabel]; found {
			deployerPods = append(deployerPods, pod)
		} else {
			d.Pods = append(d.Pods, pod)
		}
	}
	return
}

// ExtractContainersFromPods extracts all containers from the pods stored in this Data instance,
// and stores them in this Data instance.
func (d *Data) ExtractContainersFromPods() error {
	if d.Pods == nil {
		return fmt.Errorf("Empty pods")
	}
	d.Containers = []kapi.Container{}
	for _, pod := range d.Pods {
		d.Containers = append(d.Containers, pod.Spec.Containers...)
	}
	return nil
}

func (d *Data) SetContainers(containers []interface{}) error {
	d.Containers = []kapi.Container{}
	for _, obj := range containers {
		c, ok := obj.(kapi.Container)
		if !ok {
			return fmt.Errorf("Wrong type %T for containers!", containers)
		}
		d.Containers = append(d.Containers, c)
	}
	return nil
}

func (d *Data) SetImageStreams(imageStreams []interface{}) error {
	d.ImageStreams = []imageapi.ImageStream{}
	for _, obj := range imageStreams {
		is, ok := obj.(imageapi.ImageStream)
		if !ok {
			return fmt.Errorf("Wrong type %T for imageStreams!", imageStreams)
		}
		d.ImageStreams = append(d.ImageStreams, is)
	}
	return nil
}

func (d *Data) SetBuildConfigs(buildConfigs []interface{}) error {
	d.BuildConfigs = []buildapi.BuildConfig{}
	for _, obj := range buildConfigs {
		bc, ok := obj.(buildapi.BuildConfig)
		if !ok {
			return fmt.Errorf("Wrong type %T for buildConfigs!", buildConfigs)
		}
		d.BuildConfigs = append(d.BuildConfigs, bc)
	}
	return nil
}

func (d *Data) SetBuilds(builds []interface{}) error {
	d.Builds = []buildapi.Build{}
	for _, obj := range builds {
		build, ok := obj.(buildapi.Build)
		if !ok {
			return fmt.Errorf("Wrong type %T for builds!", builds)
		}
		d.Builds = append(d.Builds, build)
	}
	return nil
}

func (d *Data) SetDeploymentConfigs(deploymentConfigs []interface{}) error {
	d.DeploymentConfigs = []deployapi.DeploymentConfig{}
	for _, obj := range deploymentConfigs {
		dc, ok := obj.(deployapi.DeploymentConfig)
		if !ok {
			return fmt.Errorf("Wrong type %T for deploymentConfigs!", deploymentConfigs)
		}
		d.DeploymentConfigs = append(d.DeploymentConfigs, dc)
	}
	return nil
}

func (d *Data) SetReplicationControllers(replicationControllers []interface{}) error {
	d.ReplicationControllers = []kapi.ReplicationController{}
	for _, obj := range replicationControllers {
		rc, ok := obj.(kapi.ReplicationController)
		if !ok {
			return fmt.Errorf("Wrong type %T for replicationControllers!", replicationControllers)
		}
		d.ReplicationControllers = append(d.ReplicationControllers, rc)
	}
	return nil
}

func (d *Data) SetEvents(events []interface{}) error {
	d.Events = []kapi.Event{}
	for _, obj := range events {
		event, ok := obj.(kapi.Event)
		if !ok {
			return fmt.Errorf("Wrong type %T for events!", events)
		}
		d.Events = append(d.Events, event)
	}
	return nil
}

// ExtractApplicationsFromDeploymentConfigs extracts all applications from the deploymentConfigs stored in this Data instance,
// and stores them in this Data instance.
// The applications are in facts labels values for the label key "application"
// We inspect the DC to retrieve the different labels values, because it is the object type that all applications are likely to have.
func (d *Data) ExtractApplicationsFromDeploymentConfigs() {
	applications := make(map[string]interface{})
	for _, dc := range d.DeploymentConfigs {
		if appName, labelExists := dc.Labels[ApplicationNameLabel]; labelExists {
			if _, exists := applications[appName]; !exists {
				applications[appName] = nil
			}
		}
	}

	d.Applications = []Application{}
	for appName := range applications {
		d.Applications = append(d.Applications, Application(appName))
	}
	sort.Sort(Applications(d.Applications))
}
