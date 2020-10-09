package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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
	if os.Getenv("SITEROOT") == "" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("loading .env: %+v", err)
		}
		if os.Getenv("SITEROOT") == "" {
			return nil, errors.New("SITEROOT envvar still unset after loading .env")
		}
	}

	return &Server{
		Port:   8080,
		root:   os.Getenv("SITEROOT"),
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

// FullPath joins the specified path elements to the server's fileroot.
func (s *Server) FullPath(elem ...string) string {
	return filepath.Join(append([]string{s.root}, elem...)...)
}

// FullPathExists checks whether a fully qualified path is present on the server's filesystem.
func (s *Server) FullPathExists(fp string) (bool, error) {
	_, err := os.Stat(fp)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RelPathExists checks whether a relative path is present on the local filesystem.
func (s *Server) RelPathExists(rel string) (bool, error) {
	return s.FullPathExists(s.FullPath(rel))
}

// TemplateFile returns the full path to the named HTML template.
func (s *Server) TemplateFile(name string) string {
	basename, ok := s.templateFiles[name]
	if !ok {
		panicf("invalid template name: %s", name)
	}
	fp := s.FullPath("static", "html", basename)
	if exists, err := s.FullPathExists(fp); err != nil {
		panic(err)
	} else if !exists {
		panic(fmt.Sprintf("HTML template file %q not found at path %q", name, fp))
	}
	return fp
}

// BlogRoot returns the full path to the directory containing all blog files.
func (s *Server) BlogRoot() string {
	return s.FullPath("static", "blog")
}
