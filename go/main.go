package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

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
	server, err := NewServer()
	if err != nil {
		return fmt.Errorf("initializing server: %+v", err)
	}
	if err := server.LoadData(); err != nil {
		return fmt.Errorf("loading server data: %+v", err)
	}
	server.SetupRoutes()
	server.Router.ServeFiles("/css/*filepath", http.Dir(server.FullPath("static/css")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.Port), server))
	return nil
}
