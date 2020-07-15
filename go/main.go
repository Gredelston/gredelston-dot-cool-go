package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("loading .env: %+v", err)
	}
	server := NewServer()
	server.SetupRoutes()
	server.Router.ServeFiles("/css/*filepath", http.Dir(server.FullPath("static/css")))

	log.Fatal(http.ListenAndServe(":8080", server))
	return nil
}
