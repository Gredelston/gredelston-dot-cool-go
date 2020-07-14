package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

// exists checks whether path is present on the local filesystem.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

// panicf formats s, and panics.
func panicf(s string, i ...interface{}) { panic(fmt.Sprintf(s, i...)) }

// Main function, handles routing.
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/about", About)
	router.GET("/blog/:slug", Blog)

	router.ServeFiles("/css/*filepath", http.Dir("static/css"))

	router.PanicHandler = RenderError

	log.Fatal(http.ListenAndServe(":8080", router))
	return nil
}
