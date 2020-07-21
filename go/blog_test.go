package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func TestBlogRootExists(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	blogRoot := s.BlogRoot()
	if exists, err := s.FullPathExists(blogRoot); err != nil {
		t.Fatal(err)
	} else if !exists {
		t.Errorf("Blog root %q not present on server", blogRoot)
	}
}

func TestAllDirsWithin(t *testing.T) {
	// Create parent dir
	tempParentDir, err := ioutil.TempDir("/tmp", "temp-parent-dir")
	if err != nil {
		t.Fatal("Creating temporary parent dir: ", err)
	}
	defer os.RemoveAll(tempParentDir)

	// Populate with the following structure:
	// temp-parent-dir/
	// | temp-child-dir-1/
	// | | temp-grandchild-dir/
	// | temp-child-dir-2/
	// | temp-child-file
	tempChildDir1, err := ioutil.TempDir(tempParentDir, "temp-child-dir-1")
	if err != nil {
		t.Fatal("Creating temporary child dir 1: ", err)
	}
	tempChildDir2, err := ioutil.TempDir(tempParentDir, "temp-child-dir-2")
	if err != nil {
		t.Fatal("Creating temporary child dir 2: ", err)
	}
	if _, err = ioutil.TempDir(tempChildDir1, "temp-grandchild-dir"); err != nil {
		t.Fatal("Creating temporary grandchild dir: ", err)
	}
	if _, err = ioutil.TempFile(tempParentDir, "temp-child-file"); err != nil {
		t.Fatal("Creating temporary child file: ", err)
	}

	childDirs, err := allDirsWithin(tempParentDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(childDirs) != 2 {
		t.Errorf("Expected 2 child dirs, but got %d: %+v", len(childDirs), childDirs)
	}
	expectedChildDirs := []string{tempChildDir1, tempChildDir2}
	expectedChildDirsReverse := []string{tempChildDir2, tempChildDir1}
	if !reflect.DeepEqual(childDirs, expectedChildDirs) && !reflect.DeepEqual(childDirs, expectedChildDirsReverse) {
		t.Errorf("Unexpected childDirs: got %+v; want %+v", childDirs, expectedChildDirs)
	}
}

// TestAllBlogPostsHaveIndexHTML ensures that each dir in BlogRoot contains index.html.
func TestAllBlogPostsHaveIndexHTML(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	blogRoot := s.BlogRoot()
	blogDirs, err := allDirsWithin(blogRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range blogDirs {
		indexPath := filepath.Join(dir, "index.html")
		if exists, err := s.FullPathExists(indexPath); err != nil {
			t.Fatalf("Failed to determine whether blog post %s has an index.html: %+v", filepath.Base(dir), err)
		} else if !exists {
			t.Errorf("Blog post %s has no index.html", filepath.Base(dir))
		}
	}
}

// TestAllBlogPostsHaveMetaJSON ensures that each dir in BlogRoot contains meta.json.
func TestAllBlogPostsHaveMetaJSON(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	blogRoot := s.BlogRoot()
	blogDirs, err := allDirsWithin(blogRoot)
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range blogDirs {
		metaPath := filepath.Join(dir, "meta.json")
		if exists, err := s.FullPathExists(metaPath); err != nil {
			t.Fatalf("Failed to determine whether blog post %s has an meta.json: %+v", filepath.Base(dir), err)
		} else if !exists {
			t.Errorf("Blog post %s has no meta.json", filepath.Base(dir))
		}
	}
}

// TestInternalLinks ensures that all INTERNAL links within blog posts yield 200 responses.
func TestInternalLinks(t *testing.T) {
	srv, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	posts, err := LoadBlogPosts(srv.BlogRoot())
	if err != nil {
		t.Fatal(err)
	}
	srv.BlogPosts = posts
	srv.Router.GET("/blog/:slug", srv.HandleBlogPost)
	hrefRe := regexp.MustCompile(`<a [^>]*href=['"]([^'"]+)['"][^>]*>`)
	var url string
	for _, post := range posts {
		submatches := hrefRe.FindAllStringSubmatch(string(post.Body), -1)
		for _, submatch := range submatches {
			url = submatch[1]
			if url[0] != '/' {
				continue
			}
			url = fmt.Sprintf("http://localhost:%d%s", srv.Port, url)
			r := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, r)
			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("blog post %s: internal url %q: got status code %d; want %d", post.Slug, url, w.Result().StatusCode, http.StatusOK)
			}
		}
	}
}

// TestImages ensures that all <img> sources within blog posts yield 200 responses.
func TestImages(t *testing.T) {
	srv, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	posts, err := LoadBlogPosts(srv.BlogRoot())
	if err != nil {
		t.Fatal(err)
	}
	srv.BlogPosts = posts
	srv.SetupRoutes()
	srcRe := regexp.MustCompile(`<img [^>]*src=['"]([^'"]+)['"][^>]*>`)
	var url string
	for _, post := range posts {
		submatches := srcRe.FindAllStringSubmatch(string(post.Body), -1)
		for _, submatch := range submatches {
			url = submatch[1]
			if url[0] != '/' {
				t.Errorf("blog post %s: image url %q: externally hosted, which is dangerous", post.Slug, url)
				continue
			}
			url = fmt.Sprintf("http://localhost:%d%s", srv.Port, url)
			r := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, r)
			if w.Result().StatusCode != http.StatusOK {
				t.Errorf("blog post %s: image url %q: got status code %d; want %d", post.Slug, url, w.Result().StatusCode, http.StatusOK)
			}
		}
	}
}

// TestExternalLinks ensures that all EXTERNAL links within blog posts yield 200 responses.
func TestExternalLinks(t *testing.T) {
	srv, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	posts, err := LoadBlogPosts(srv.BlogRoot())
	if err != nil {
		t.Fatal(err)
	}
	hrefRe := regexp.MustCompile(`<a [^>]*href=['"]([^'"]+)['"][^>]*>`)
	var url string
	for _, post := range posts {
		submatches := hrefRe.FindAllStringSubmatch(string(post.Body), -1)
		for _, submatch := range submatches {
			url = submatch[1]
			if url[0] == '/' {
				continue
			}
			// First try the fast way, which does not require HTTP headers
			resp, err := http.Get(url)
			if err != nil {
				t.Errorf("blog post %s: external url %q: %+v", post.Slug, url, err)
			}
			if resp.StatusCode == http.StatusOK {
				continue
			}
			// If that failed, try mimicking browser headers
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Errorf("blog post %s: url %q: creating request: %+v", post.Slug, url, err)
			}
			req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
			if resp, err = http.DefaultClient.Do(req); err != nil {
				t.Errorf("blog post %s: external url %q: %+v", post.Slug, url, err)
			} else if resp.StatusCode != http.StatusOK {
				t.Errorf("blog post %s: external url %q: got status code %d; want %d", post.Slug, url, resp.StatusCode, http.StatusOK)
			}
		}
	}
}
