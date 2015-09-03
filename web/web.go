// Package web provides an HTTP server for the dashboard application
package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/tylerb/graceful"
)

// RunHttpServer runs an HTTP server on a port defined by the PORT env var (default to 8080)
func RunHttpServer() {
	port := Getenv("PORT", "8080")
	publicDir := Getenv("PUBLIC_DIR", "public")

	c := NewContext()

	router := httprouter.New()
	router.GET("/", c.Home)

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir(publicDir)),
	)

	n.UseHandler(router)

	log.Printf("Starting openshift-dashboard on port %v\n", port)
	graceful.Run(":"+port, 10*time.Second, n)
}

// Getenv returns the value of the env var with the given name,
// or fallback to the given default value.
func Getenv(envVarName string, defaultValue string) string {
	if value := os.Getenv(envVarName); len(value) != 0 {
		return value
	}
	return defaultValue
}
