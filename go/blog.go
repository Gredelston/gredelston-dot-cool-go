package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// BlogPost contains the data of a single blog post.
type BlogPost struct {
	Body       template.HTML
	Hidden     bool
	PostedDate time.Time
	rawDate    string `json:Date`
	Slug       string
	Title      string
	Tags       []string
}

// allDirsWithin returns an array of directories within some dir parent.
func allDirsWithin(parent string) ([]string, error) {
	var childDirs []string
	children, err := ioutil.ReadDir(parent)
	if err != nil {
		return nil, fmt.Errorf("reading dir %s: %+v", parent, err)
	}
	for _, child := range children {
		if child.IsDir() {
			childDirs = append(childDirs, filepath.Join(parent, child.Name()))
		}
	}
	return childDirs, nil
}

// loadBlogPost reads a blog-post directory and returns a populated BlogPost object.
func loadBlogPost(dir string) (*BlogPost, error) {
	// Read index.md
	indexFP := filepath.Join(dir, "index.md")
	b, err := ioutil.ReadFile(indexFP)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", indexFP, err)
	}
	body := string(b)
	body = strings.TrimSpace(body)
	body = strings.ReplaceAll(body, "\n", "<br>")
	body = regexp.MustCompile(`\[((?:[^\]\r\n]|\\\])+)\]\((\S+)\)`).ReplaceAllString(body, `<a href=${2}>${1}</a>`)

	// Read meta.json
	metaFP := filepath.Join(dir, "meta.json")
	b, err = ioutil.ReadFile(metaFP)
	if err != nil {
		return nil, fmt.Errorf("readings %s: %+v", metaFP, err)
	}
	post := &BlogPost{}
	if err = json.Unmarshal(b, post); err != nil {
		return nil, fmt.Errorf("unmarshaling %s: %+v", metaFP, err)
	}

	// Populate BlogPost and return
	post.PostedDate = time.Now()
	post.Body = template.HTML(body)
	post.Slug = filepath.Base(dir)
	return post, nil
}

// LoadBlogPosts converts each directory within blogRoot into a BlogPost object.
func LoadBlogPosts(blogRoot string) ([]*BlogPost, error) {
	blogPostDirs, err := allDirsWithin(blogRoot)
	if err != nil {
		return nil ,err
	}
	posts := make([]*BlogPost, len(blogPostDirs))
	for i, dir := range blogPostDirs {
		if post, err := loadBlogPost(dir); err != nil {
			return nil, fmt.Errorf("loading blog post %s: %+v", filepath.Base(dir), err)
		} else {
			posts[i] = post
		}
	}
	return posts, nil
}

// BlogPostBySlug searches the server's loaded BlogPosts for one matching by Slug.
func (s *Server) BlogPostBySlug(slug string) (bp *BlogPost, ok bool) {
	for _, bp := range s.BlogPosts {
		if bp.Slug == slug {
			return bp, true
		}
	}
	return nil, false
}
