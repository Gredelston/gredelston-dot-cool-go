package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gredelston/gredelston-dot-cool-go/server"
)

// Main function, handles routing.
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s, err := server.NewServer()
	if err != nil {
		return fmt.Errorf("initializing server: %+v", err)
	}
	if err := s.LoadData(); err != nil {
		return fmt.Errorf("loading server data: %+v", err)
	}
	s.SetupRoutes()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s))
	return nil
}
