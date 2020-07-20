package main

import (
	"fmt"
	"html/template"
	"net/http"
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
	BlogData *BlogPost

	// BlogPosts contains the list of blog posts to display on blogroll.html.
	BlogPosts []*BlogPost
}

// RenderPage generically render a specified page template with custom data.
func (s *Server) RenderPage(name string, w http.ResponseWriter, data PageData) {
	if len(data.Navs) == 0 {
		data.Navs = NewNavs()
	}
	mainFP := s.TemplateFile(name)
	t := template.New(path.Base(mainFP))
	t = template.Must(t.ParseFiles(mainFP, s.TemplateFile("header"), s.TemplateFile("footer"), s.TemplateFile("navbar"), s.TemplateFile("blogroll")))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

// RenderError renders generic error text to the browser.
func (s *Server) RenderError(w http.ResponseWriter, _ *http.Request, err interface{}) {
	s.RenderPage("text", w, PageData{
		Title: "Arrow'd!",
		Body:  fmt.Sprintf("Internal server error: %+v", err),
	})
}
