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
			"now":      "now.html",
		},
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// LoadData initializes any data that the server will depend on.
func (s *Server) LoadData() error {
	// Blog posts
	br, err := s.BlogRoot()
	if err != nil {
		return err
	}
	posts, err := LoadBlogPosts(br)
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
func (s *Server) TemplateFile(name string) (string, error) {
	basename, ok := s.templateFiles[name]
	if !ok {
		return "", fmt.Errorf("invalid template name: %s", name)
	}
	htmlFP, err := utils.Asset("html", basename)
	if err != nil {
		return "", fmt.Errorf("determining asset path for html/%s: %+v", basename, err)
	}
	if exists, err := utils.PathExists(htmlFP); err != nil {
		return "", err
	} else if !exists {
		return "", fmt.Errorf("HTML template file %q not found at path %q", name, htmlFP)
	}
	return htmlFP, nil
}

// TemplateMustExist returns the full path to the named HTML template, or panics if it doesn't exist.
func (s *Server) TemplateMustExist(name string) string {
	fp, err := s.TemplateFile(name)
	if err != nil {
		panic(err)
	}
	return fp
}

// BlogRoot returns the full path to the directory containing all blog files.
func (s *Server) BlogRoot() (string, error) {
	br, err := utils.Asset("blog")
	if err != nil {
		return "", fmt.Errorf("determining blog root: %+v", err)
	}
	return br, nil
}
