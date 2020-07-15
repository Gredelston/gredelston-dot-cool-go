package main

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) SetupRoutes() {
	s.Router.GET("/", s.HandleIndex)
	s.Router.GET("/about", s.HandleAbout)
	s.Router.GET("/blog/:slug", s.HandleBlogPost)
	s.Router.PanicHandler = s.RenderError
}

// HandleIndex defines a GET method to handle / routes.
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RenderPage("index", w, PageData{
		Title: "Home",
		ExtraStylesheets: []string{("/css/index.css")},
		Navs: NewNavsWithActive(NavHome),
	})
}

// HandleAbout defines a GET method to handle /about routes.
func (s *Server) HandleAbout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RenderPage("text", w, PageData{
		Title: "About Me",
		Body: `My name is Greg. I'm a software developer in Boulder, Colorado.`,
		Keywords: []string{"UnderConstruction"},
		Navs: NewNavsWithActive(NavAbout),
	})
}

// HandleBlogPost defines a GET method to handle /blog/:slug routes.
func (s *Server) HandleBlogPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("slug")
	blogDir := s.FullPath("static/blog", slug)
	if blogDirExists, err := s.FullPathExists(blogDir); err != nil {
		panicf("finding blogdir for slug %q: %+v", slug, err)
	} else if !blogDirExists {
		panicf("did not find blogdir for slug %q", slug)
	}
	blogFile := filepath.Join(blogDir, "index.md")
	if blogFileExists, err := s.FullPathExists(blogFile); err != nil {
		panicf("Finding blogfile for slug %q: %+v", slug, err)
	} else if !blogFileExists {
		panicf("Did not find blogfile for slug %q", slug)
	}
	title := slug
	date := time.Now()
	tags := []string{"CoolTag1", "UncoolTagLambda"}
	body := "Wow, what a cool blog post!"
	s.RenderPage("blog", w, PageData{
		Title: title,
		BlogData: BlogData{
			Title: title,
			Date: date,
			Tags: tags,
			Body: body,
		},
	})
}
