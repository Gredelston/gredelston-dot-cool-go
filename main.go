package main

import (
	"fmt"
	"html/template"
	"net/http"
	"log"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseFiles("template/page.html")
	if err != nil {
		fmt.Errorf("Failed to parse files: %w", err)
		return
	}
	data := map[string]string {"body": "Welcome to my website!"}
	if err = t.Execute(w, data); err != nil {
		fmt.Errorf("Failed to execute template: ", err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}
