package server

import (
	"fmt"
	"net/http"

	"github.com/gredelston/gredelston-dot-cool-go/utils"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	BlogPosts     []*BlogPost
	Port          int
	root          string
	Router        *httprouter.Router
	templateFiles map[string]string
}

func NewServer() (*Server, error) {
	return &Server{
		Port:   8080,
		Router: httprouter.New(),
		templateFiles: map[string]string{
			// Generic template files
			"header": "header.html",
			"footer": "footer.html",
			"navbar": "navbar.html",

			// Page-specific template files
			"about":    "about.html",
			"blog":     "blog.html",
			"blogroll": "blogroll.html",
			"index":    "index.html",
			"text":     "text.html",
		},
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// LoadData initializes any data that the server will depend on.
func (s *Server) LoadData() error {
	// Blog posts
	posts, err := LoadBlogPosts(s.BlogRoot())
	if err != nil {
		return fmt.Errorf("loading blog data: %+v", err)
	}
	s.BlogPosts = posts
	for _, p := range posts {
		s.Router.ServeFiles(fmt.Sprintf("/assets/%s/*filepath", p.Slug), http.Dir(p.ServerPath))
	}

	return nil
}

// TemplateFile returns the full path to the named HTML template.
func (s *Server) TemplateFile(name string) string {
	basename, ok := s.templateFiles[name]
	if !ok {
		utils.Panicf("invalid template name: %s", name)
	}
	fp := utils.StaticPath("html", basename)
	if exists, err := utils.PathExists(fp); err != nil {
		panic(err)
	} else if !exists {
		panic(fmt.Sprintf("HTML template file %q not found at path %q", name, fp))
	}
	return fp
}

// BlogRoot returns the full path to the directory containing all blog files.
func (s *Server) BlogRoot() string {
	return utils.StaticPath("blog")
}
