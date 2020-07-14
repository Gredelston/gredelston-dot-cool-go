package main

import (
	"fmt"
	"net/http"
	"html/template"
	"time"
)

const (
	// Generic template filepaths
	headerTFP = "template/header.html"
	footerTFP = "template/footer.html"
	navbarTFP = "template/navbar.html"

	// Page-specific template filepaths
	aboutTFP = "template/about.html"
	blogTFP  = "template/blog.html"
	indexTFP = "template/index.html"
	textTFP  = "template/text.html"
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
	BlogData BlogData
}

// BlogData contains the data needed to render a blog post.
type BlogData struct {
	Title string
	Date  time.Time
	Tags  []string
	Body  string
}

// RenderPage generically render a specified template (fp) with custom data.
func RenderPage(fp string, w http.ResponseWriter, data PageData) {
	if len(data.Navs) == 0 {
		data.Navs = NewNavs()
	}
	t := template.Must(template.ParseFiles(fp, headerTFP, footerTFP, navbarTFP))
	if err := t.Execute(w, data); err != nil { panic(err) }
}

// RenderError renders generic error text to the browser.
func RenderError(w http.ResponseWriter, _ *http.Request, err interface{}) {
	RenderPage(textTFP, w, PageData{
		Title: "Arrow'd!",
		Body: fmt.Sprintf("Internal server error: %+v", err),
	})
}
