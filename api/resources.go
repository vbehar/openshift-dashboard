package api

// ResourceType describes the possible types of resources.
type ResourceType string

const (
	// ResourceTypeApplication represents an application, which is just a label on objects
	ResourceTypeApplication           ResourceType = "application"
	ResourceTypeProject               ResourceType = "project"
	ResourceTypeRoute                 ResourceType = "route"
	ResourceTypeService               ResourceType = "service"
	ResourceTypePod                   ResourceType = "pod"
	ResourceTypeContainer             ResourceType = "container"
	ResourceTypeImageStream           ResourceType = "imagestream"
	ResourceTypeBuildConfig           ResourceType = "buildconfig"
	ResourceTypeBuild                 ResourceType = "build"
	ResourceTypeDeploymentConfig      ResourceType = "deploymentconfig"
	ResourceTypeReplicationController ResourceType = "replicationcontroller"
	ResourceTypeEvent                 ResourceType = "event"
)

var (
	ResourceTypeAll []ResourceType = []ResourceType{
		ResourceTypeApplication,
		ResourceTypeProject,
		ResourceTypeRoute,
		ResourceTypeService,
		ResourceTypePod,
		ResourceTypeContainer,
		ResourceTypeImageStream,
		ResourceTypeBuildConfig,
		ResourceTypeBuild,
		ResourceTypeDeploymentConfig,
		ResourceTypeReplicationController,
		ResourceTypeEvent,
	}
)
