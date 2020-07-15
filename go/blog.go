package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

// BlogPost contains the data of a single blog post.
type BlogPost struct {
	Body   string
	Date   time.Time
	Hidden bool
	Slug   string
	Title  string
	Tags   []string
}

func (s *Server) LoadBlogData() error {
	blogRoot := s.FullPath("static/blog")
	if exists, err := s.FullPathExists(blogRoot); err != nil {
		return fmt.Errorf("determining whether blog root exists: %+v", err)
	} else if !exists {
		return fmt.Errorf("blog root %s does not exist", blogRoot)
	}
	dirs, err := ioutil.ReadDir(blogRoot)
	if err != nil {
		return fmt.Errorf("failed to read files from blog root: %+v", err)
	}
	var posts []BlogPost
	for _, dir := range dirs {
		if !dir.IsDir() { continue }
		indexPath := filepath.Join(blogRoot, dir.Name(), "index.md")
		fmt.Println(indexPath)
		if exists, err := s.FullPathExists(indexPath); err != nil {
			return fmt.Errorf("determining whether blog post %s has an index.md: %+v", dir, err)
		} else if !exists {
			continue
		}
		posts = append(posts, BlogPost{
			Body: "Hello!",
			Date: time.Now(),
			Slug: dir.Name(),
			Title: "My cool title!",
		})
	}
	s.BlogPosts = posts
	return nil
}

// BlogPostBySlug searches the server's loaded BlogPosts for one matching by Slug.
func (s *Server) BlogPostBySlug(slug string) (bp *BlogPost, ok bool) {
	for _, bp := range s.BlogPosts {
		if bp.Slug == slug {
			return &bp, true
		}
	}
	return nil, false
}
