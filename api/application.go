package api

const (
	// ApplicationNameLabel is the name of the label used to store the application
	ApplicationNameLabel = "application"
)

// Application is a concept used here to group multiple objects together
// It is different from the project, because the same project can be used for multiple applications
// Or some projects can have resources (like BuildConfigs and ImageStreams) but no applications.
// Objects are linked to an application by a label.
type Application string

// Name returns the name of the application
// (just to be consistent with others objects that have a "Name" attribute)
func (app *Application) Name() string {
	return string(*app)
}

// Applications is just a slice of Application
// used for sorting applications
type Applications []Application

func (apps Applications) Len() int           { return len(apps) }
func (apps Applications) Less(i, j int) bool { return apps[i] < apps[j] }
func (apps Applications) Swap(i, j int)      { apps[i], apps[j] = apps[j], apps[i] }
