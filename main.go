package main

import (
	"fmt"
	"net/http"
	"log"

	"github.com/julienschmidt/httprouter"
)

// panicf formats s, and panics.
func panicf(s string, i ...interface{}) { panic(fmt.Sprintf(s, i...)) }

// Main function, handles routing.
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/about", About)
	router.GET("/blog/:slug", Blog)

	router.ServeFiles("/css/*filepath", http.Dir("static/css"))

	router.PanicHandler = RenderError

	log.Fatal(http.ListenAndServe(":8080", router))
}
