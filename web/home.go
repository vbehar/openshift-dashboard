package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vbehar/openshift-dashboard/api"

	"github.com/julienschmidt/httprouter"
)

// Data represents the data retrieved from the API, and exposed to view
type Data struct {
	*api.Data
}

// Title returns the title of the page
// using either the env var DASHBOARD_TITLE
// or a default title.
func (d *Data) Title() string {
	if title := os.Getenv("DASHBOARD_TITLE"); len(title) > 0 {
		return title
	}

	// default value
	return "openshift-dashboard"
}

// Home answers HTTP requests by loading data for all resource types and using the "home" view
func (c *Context) Home(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	d, err := c.ClientWrapper.LoadData(api.ResourceTypeAll...)
	if err != nil {
		fmt.Fprintf(w, "failed to load data: %v", err)
		return
	}

	data := &Data{d}

	c.Render.HTML(w, http.StatusOK, "home", data)
}
