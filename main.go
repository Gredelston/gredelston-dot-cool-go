package main

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
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

// PageData contains the data needed by renderPage.
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

// renderPage generically render a specified template (fp) with custom data.
func renderPage(fp string, w http.ResponseWriter, data PageData) {
	if len(data.Navs) == 0 {
		data.Navs = NewNavs()
	}
	t := template.Must(template.ParseFiles(fp, headerTFP, footerTFP, navbarTFP))
	if err := t.Execute(w, data); err != nil { panic(err) }
}

// renderError renders generic error text to the browser.
func renderError(err error, w http.ResponseWriter) {
	renderPage(textTFP, w, PageData{
		Title: "Arrow'd!",
		Body: fmt.Sprintf("%+v", err),
	})
}

// Index defines a GET method to route homepage.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderPage(indexTFP, w, PageData{
		Title: "Home",
		ExtraStylesheets: []string{"/css/index.css"},
		Navs: NewNavsWithActive(NavHome),
	})
}

// About defines a GET method to route "About Me" page.
func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderPage(textTFP, w, PageData{
		Title: "About Me",
		Body: `My name is Greg. I'm a software developer in Boulder, Colorado.`,
		Keywords: []string{"UnderConstruction"},
		Navs: NewNavsWithActive(NavAbout),
	})
}

// exists checks whether path is present on the local filesystem.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

// GET method to route blog pages.
func Blog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("slug")
	blogDir := filepath.Join("static", "blog", slug)
	if blogDirExists, err := exists(blogDir); err != nil {
		renderError(fmt.Errorf("Finding blogdir for slug %q: %+v", slug, err), w)
		return
	} else if !blogDirExists {
		renderError(fmt.Errorf("Did not find blogdir for slug %q", slug), w)
		return
	}
	blogFile := filepath.Join(blogDir, "index.md")
	if blogFileExists, err := exists(blogFile); err != nil {
		renderError(fmt.Errorf("Finding blogfile for slug %q: %+v", slug, err), w)
		return
	} else if !blogFileExists {
		renderError(fmt.Errorf("Did not find blogfile for slug %q", slug), w)
		return
	}
	title := slug
	date := time.Now()
	tags := []string{"CoolTag1", "UncoolTagLambda"}
	body := "Wow, what a cool blog post!"
	renderPage(blogTFP, w, PageData{
		Title: title,
		BlogData: BlogData{
			Title: title,
			Date: date,
			Tags: tags,
			Body: body,
		},
	})
}

// Main function, handles routing.
func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/about", About)
	router.GET("/blog/:slug", Blog)

	router.ServeFiles("/css/*filepath", http.Dir("static/css"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
