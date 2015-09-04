package web

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// StatsHandler answers HTTP requests with the stats in JSON format
func (c *Context) StatsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	s, err := json.Marshal(c.Stats.Data())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(s)
}
