package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	BlogPosts []BlogPost
	root      string
	Router    *httprouter.Router
	templates map[string]string
}

func NewServer() (*Server, error) {
	if os.Getenv("SITEROOT") == "" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("loading .env: %+v", err)
		}
	}

	return &Server{
		root:      os.Getenv("SITEROOT"),
		Router:    httprouter.New(),
		templates: map[string]string{
			// Generic template filepaths
			"header": "template/header.html",
			"footer": "template/footer.html",
			"navbar": "template/navbar.html",

			// Page-specific template filepaths
			"about": "template/about.html",
			"blog":  "template/blog.html",
			"index": "template/index.html",
			"text":  "template/text.html",
		},
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// LoadData initializes any data that the server will depend on.
func (s *Server) LoadData() error {
	if err := s.LoadBlogPosts(); err != nil {
		return fmt.Errorf("loading blog data: %+v", err)
	}
	return nil
}

// FullPath joins the specified path elements to the server's fileroot.
func (s *Server) FullPath(elem ...string) string {
	return filepath.Join(append([]string{s.root}, elem...)...)
}

// FullPathExists checks whether a fully qualified path is present on the server's filesystem.
func (s *Server) FullPathExists(fp string) (bool, error) {
	_, err := os.Stat(fp)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

// RelPathExists checks whether a relative path is present on the local filesystem.
func (s *Server) RelPathExists(rel string) (bool, error) {
	return s.FullPathExists(s.FullPath(rel))
}

// Template returns the full path to the named HTML template.
func (s *Server) Template(name string) string {
	rel, ok := s.templates[name]
	if !ok {
		panicf("invalid template name: %s", name)
	}
	fp := s.FullPath(rel)
	if exists, err := s.FullPathExists(fp); err != nil {
		panic(err)
	} else if !exists {
		panic(fmt.Sprintf("HTML template file %q not found", name))
	}
	return fp
}

// BlogRoot returns the full path to the blog root.
func (s *Server) BlogRoot() string {
	return s.FullPath("static", "blog")
}
