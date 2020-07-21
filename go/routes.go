package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (s *Server) SetupRoutes() {
	s.Router.GET("/", s.HandleIndex)
	s.Router.GET("/about", s.HandleAbout)
	s.Router.GET("/blog/:slug", s.HandleBlogPost)
	s.Router.ServeFiles("/css/*filepath", http.Dir(s.FullPath("static/css")))
	s.Router.ServeFiles("/images/*filepath", http.Dir(s.FullPath("static/images")))
	s.Router.PanicHandler = s.RenderError
}

// HandleIndex defines a GET method to handle / routes.
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RenderPage("index", w, PageData{
		Title:            "Home",
		ExtraStylesheets: []string{("/css/index.css")},
		Navs:             NewNavsWithActive(NavHome),
		BlogPosts:        s.BlogPosts,
	})
}

// HandleAbout defines a GET method to handle /about routes.
func (s *Server) HandleAbout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RenderPage("text", w, PageData{
		Title:    "About Me",
		Body:     `My name is Greg. I'm a software developer in Boulder, Colorado.`,
		Keywords: []string{"UnderConstruction"},
		Navs:     NewNavsWithActive(NavAbout),
	})
}

// HandleBlogPost defines a GET method to handle /blog/:slug routes.
func (s *Server) HandleBlogPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("slug")
	bp, ok := s.BlogPostBySlug(slug)
	if !ok {
		panicf("no blog post found for slug %s", slug)
	}
	title := bp.Slug
	s.RenderPage("blog", w, PageData{
		Title:    title,
		BlogData: bp,
	})
}
