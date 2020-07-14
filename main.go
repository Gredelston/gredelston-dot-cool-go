package main

import (
	"html/template"
	"net/http"
	"log"
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

// Generic function to render a specified template (fp) with custom data.
func renderPage(fp string, w http.ResponseWriter, data PageData) {
	if len(data.Navs) == 0 {
		data.Navs = NewNavs()
	}
	t := template.Must(template.ParseFiles(fp, headerTFP, footerTFP, navbarTFP))
	if err := t.Execute(w, data); err != nil { panic(err) }
}

// GET method to route homepage.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderPage(indexTFP, w, PageData{
		Title: "Home",
		ExtraStylesheets: []string{"/css/index.css"},
		Navs: NewNavsWithActive(NavHome),
	})
}

// GET method to route "About Me" page.
func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	renderPage(textTFP, w, PageData{
		Title: "About Me",
		Body: `My name is Greg. I'm a software developer in Boulder, Colorado.`,
		Keywords: []string{"UnderConstruction"},
		Navs: NewNavsWithActive(NavAbout),
	})
}

// GET method to route blog pages.
func Blog(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	title := "My Fake Blog Post"
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
