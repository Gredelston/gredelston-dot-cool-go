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

func (s *Server) LoadBlogData() error {
	blogPostDirs, err := allDirsWithin(s.BlogRoot())
	if err != nil { return err }
	var posts []BlogPost
	for _, dir := range blogPostDirs {
		indexPath := filepath.Join(dir, "index.md")
		fmt.Println(indexPath)
		if exists, err := s.FullPathExists(indexPath); err != nil {
			return fmt.Errorf("determining whether blog post %s has an index.md: %+v", dir, err)
		} else if !exists {
			continue
		}
		posts = append(posts, BlogPost{
			Body: "Hello!",
			Date: time.Now(),
			Slug: filepath.Base(dir),
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
