package api

import (
	"fmt"

	buildapi "github.com/openshift/origin/pkg/build/api"
	deployapi "github.com/openshift/origin/pkg/deploy/api"
	imageapi "github.com/openshift/origin/pkg/image/api"
	routeapi "github.com/openshift/origin/pkg/route/api"

	kapi "k8s.io/kubernetes/pkg/api"
)

// FilterByApplication returns all objects that belongs to the given application
func FilterByApplication(objects interface{}, application string) ([]interface{}, error) {
	return FilterByLabelValue(objects, ApplicationNameLabel, application)
}

// FilterByLabelValue returns all objects that have the given key/value label
func FilterByLabelValue(objects interface{}, labelKey string, labelValue string) ([]interface{}, error) {
	results := []interface{}{}

	switch objects.(type) {

	case []routeapi.Route:
		for _, route := range objects.([]routeapi.Route) {
			if hasLabelValue(route.ObjectMeta, labelKey, labelValue) {
				results = append(results, route)
			}
		}

	case []kapi.Service:
		for _, service := range objects.([]kapi.Service) {
			if hasLabelValue(service.ObjectMeta, labelKey, labelValue) {
				results = append(results, service)
			}
		}

	case []imageapi.ImageStream:
		for _, is := range objects.([]imageapi.ImageStream) {
			if hasLabelValue(is.ObjectMeta, labelKey, labelValue) {
				results = append(results, is)
			}
		}

	case []buildapi.BuildConfig:
		for _, bc := range objects.([]buildapi.BuildConfig) {
			if hasLabelValue(bc.ObjectMeta, labelKey, labelValue) {
				results = append(results, bc)
			}
		}

	case []deployapi.DeploymentConfig:
		for _, dc := range objects.([]deployapi.DeploymentConfig) {
			if hasLabelValue(dc.ObjectMeta, labelKey, labelValue) {
				results = append(results, dc)
			}
		}

	default:
		return nil, fmt.Errorf("Unsupported transformation of type %T", objects)
	}

	return results, nil
}

// FilterByNamespace returns all objects that belongs to the given namespace (project)
func FilterByNamespace(objects interface{}, namespace string) ([]interface{}, error) {
	results := []interface{}{}

	switch objects.(type) {

	case []routeapi.Route:
		for _, route := range objects.([]routeapi.Route) {
			if route.Namespace == namespace {
				results = append(results, route)
			}
		}

	case []kapi.Service:
		for _, service := range objects.([]kapi.Service) {
			if service.Namespace == namespace {
				results = append(results, service)
			}
		}

	case []imageapi.ImageStream:
		for _, is := range objects.([]imageapi.ImageStream) {
			if is.Namespace == namespace {
				results = append(results, is)
			}
		}

	case []buildapi.BuildConfig:
		for _, bc := range objects.([]buildapi.BuildConfig) {
			if bc.Namespace == namespace {
				results = append(results, bc)
			}
		}

	case []deployapi.DeploymentConfig:
		for _, dc := range objects.([]deployapi.DeploymentConfig) {
			if dc.Namespace == namespace {
				results = append(results, dc)
			}
		}

	default:
		return nil, fmt.Errorf("Unsupported transformation of type %T", objects)
	}

	return results, nil
}

// hasLabelValue returns true if the given object has the given key/value label
func hasLabelValue(meta kapi.ObjectMeta, labelKey string, labelValue string) bool {
	if value, found := meta.Labels[labelKey]; found && value == labelValue {
		return true
	}
	return false
}
