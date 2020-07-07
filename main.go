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
	navbarTFP = "template/navbar.html"

	// Page-specific template filepaths
	aboutTFP = "template/about.html"
	indexTFP = "template/index.html"
	textTFP = "template/text.html"
)

type pageData struct {
	Title, Body string
	ExtraStylesheets []string
}

func renderPage(fp string, w http.ResponseWriter, data pageData) {
	t := template.Must(template.ParseFiles(fp, headerTFP, footerTFP, navbarTFP))
	if err := t.Execute(w, data); err != nil {
		panic(fmt.Errorf("Executing template %q: %w", fp, err))
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderPage(indexTFP, w, pageData{
		Title: "Home",
		ExtraStylesheets: []string{"/css/index.css"},
	})
}

func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderPage(textTFP, w, pageData{
		Title: "About Me",
		Body: `My name is Greg. I'm a software developer in Boulder, Colorado.`,
	})
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Hello, %s!\n", ps.ByName("name"))
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/about", About)
	router.GET("/hello/:name", Hello)

	router.ServeFiles("/css/*filepath", http.Dir("static/css"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
