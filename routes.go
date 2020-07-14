package main

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Index defines a GET method to route homepage.
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderPage(indexTFP, w, PageData{
		Title: "Home",
		ExtraStylesheets: []string{"/css/index.css"},
		Navs: NewNavsWithActive(NavHome),
	})
}

// About defines a GET method to route "About Me" page.
func About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	RenderPage(textTFP, w, PageData{
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
		panicf("finding blogdir for slug %q: %+v", slug, err)
	} else if !blogDirExists {
		panicf("did not find blogdir for slug %q", slug)
	}
	blogFile := filepath.Join(blogDir, "index.md")
	if blogFileExists, err := exists(blogFile); err != nil {
		panicf("Finding blogfile for slug %q: %+v", slug, err)
	} else if !blogFileExists {
		panicf("Did not find blogfile for slug %q", slug)
	}
	title := slug
	date := time.Now()
	tags := []string{"CoolTag1", "UncoolTagLambda"}
	body := "Wow, what a cool blog post!"
	RenderPage(blogTFP, w, PageData{
		Title: title,
		BlogData: BlogData{
			Title: title,
			Date: date,
			Tags: tags,
			Body: body,
		},
	})
}
