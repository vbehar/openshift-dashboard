// Package web provides an HTTP server for the dashboard application
package web

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
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

	n.Run(":" + port)
}

// Getenv returns the value of the env var with the given name,
// or fallback to the given default value.
func Getenv(envVarName string, defaultValue string) string {
	if value := os.Getenv(envVarName); len(value) != 0 {
		return value
	}
	return defaultValue
}
