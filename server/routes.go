package server

import (
	"fmt"
	"net/http"

	"github.com/gredelston/gredelston-dot-cool-go/utils"
	"github.com/julienschmidt/httprouter"
)

func (s *Server) SetupRoutes() error {
	s.Router.GET("/", s.HandleIndex)
	s.Router.GET("/about", s.HandleAbout)
	s.Router.GET("/now", s.HandleNow)
	s.Router.GET("/blog/:slug", s.HandleBlogPost)
	css, err := utils.Asset("css")
	if err != nil {
		return fmt.Errorf("getting css path: %+v", err)
	}
	s.Router.ServeFiles("/css/*filepath", http.Dir(css))
	images, err := utils.Asset("images")
	if err != nil {
		return fmt.Errorf("getting images path: %+v", err)
	}
	s.Router.ServeFiles("/images/*filepath", http.Dir(images))
	s.Router.PanicHandler = s.RenderError
	return nil
}

// HandleIndex routes / to the index.html template.
func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RenderPage("index", w, PageData{
		Title:            "Home",
		ExtraStylesheets: []string{("/css/index.css")},
		Navs:             NewNavsWithActive(NavHome),
		BlogPosts:        s.BlogPosts,
	})
}

// HandleNow routes /now to the now.html template.
func (s *Server) HandleNow(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.RenderPage("now", w, PageData{
		Title:    "Now",
		Navs:     NewNavsWithActive(NavNow),
		Keywords: []string{"now"},
	})
}

// HandleAbout routes /about to a simple Text page.
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
		utils.Panicf("no blog post found for slug %s", slug)
	}
	title := bp.Slug
	s.RenderPage("blog", w, PageData{
		Title:    title,
		BlogData: bp,
	})
}
