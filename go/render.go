package main

import (
	"fmt"
	"net/http"
	"html/template"
	"path"
)

// PageData contains the data needed by RenderPage.
type PageData struct {
	// Title is the title of the page.
	// Body is used by some templates (e.g. text.html) for the main text body.
	Title, Body string

	// ExtraStylesheets are page-specific CSS files to be imported.
	// Keywords are page-specific metadata keywords.
	ExtraStylesheets, Keywords []string

	// Navs are the items to render in the navbar.
	Navs []Nav

	// BlogData contains data used for blog-type posts.
	BlogData BlogPost
}

// RenderPage generically render a specified page template with custom data.
func (s *Server) RenderPage(name string, w http.ResponseWriter, data PageData) {
	if len(data.Navs) == 0 {
		data.Navs = NewNavs()
	}
	mainFP := s.Template(name)
	t := template.New(path.Base(mainFP))
	t = template.Must(t.ParseFiles(mainFP, s.Template("header"), s.Template("footer"), s.Template("navbar")))
	if err := t.Execute(w, data); err != nil { panic(err) }
}

// RenderError renders generic error text to the browser.
func (s *Server) RenderError(w http.ResponseWriter, _ *http.Request, err interface{}) {
	s.RenderPage("text", w, PageData{
		Title: "Arrow'd!",
		Body: fmt.Sprintf("Internal server error: %+v", err),
	})
}
