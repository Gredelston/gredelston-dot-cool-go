package main

import (
	"fmt"
	"html/template"
	"net/http"
	"log"

	"github.com/julienschmidt/httprouter"
)

const (
	// Generic template filepaths
	headerTFP = "template/header.html"
	footerTFP = "template/footer.html"

	// Page-specific template filepaths
	indexTFP = "template/index.html"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.Must(template.ParseFiles(indexTFP, headerTFP, footerTFP))
	if err := t.Execute(w, map[string]string{"title": "Index"}); err != nil {
		panic(err)
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
