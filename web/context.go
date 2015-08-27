package web

import (
	"html/template"
	"os"
	"strings"

	"github.com/vbehar/openshift-dashboard/api"

	"github.com/unrolled/render"
)

// Context is a web context used to answer to requests.
type Context struct {
	ClientWrapper *api.ClientWrapper
	Render        *render.Render
}

// NewContext builds a new Context instance
func NewContext() *Context {
	r := render.New(render.Options{
		IsDevelopment: isDevEnv(),
		Funcs: []template.FuncMap{{
			"filterByNamespace":   api.FilterByNamespace,
			"filterByApplication": api.FilterByApplication,
			"filterByLabelValue":  api.FilterByLabelValue,
		}},
	})

	cacheEnabled := !isDevEnv()
	clientWrapper := api.NewClientWrapper(cacheEnabled)

	return &Context{
		ClientWrapper: clientWrapper,
		Render:        r,
	}
}

// isDevEnv returns true if we are running in "dev" env
// It checks the value of the GO_ENV env var (it should be equals to "dev")
func isDevEnv() bool {
	goEnv := os.Getenv("GO_ENV")
	if strings.ToLower(goEnv) == "dev" {
		return true
	}
	return false
}
