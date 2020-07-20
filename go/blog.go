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

// loadBlogPost reads a blog-post directory and returns a populated BlogPost object.
func loadBlogPost(dir string) (BlogPost, error) {
	return BlogPost{
		Body: "Hello!",
		Date: time.Now(),
		Slug: filepath.Base(dir),
		Title: "My cool title!",
	}, nil
}

func (s *Server) LoadBlogData() error {
	blogPostDirs, err := allDirsWithin(s.BlogRoot())
	if err != nil {
		return err
	}
	posts := make([]BlogPost, len(blogPostDirs))
	for i, dir := range blogPostDirs {
		indexPath := filepath.Join(dir, "index.md")
		fmt.Println(indexPath)
		if exists, err := s.FullPathExists(indexPath); err != nil {
			return fmt.Errorf("determining whether blog post %s has an index.md: %+v", filepath.Base(dir), err)
		} else if !exists {
			continue
		}
		if post, err := loadBlogPost(dir); err != nil {
			return fmt.Errorf("loading blog post %s: %+v", filepath.Base(dir), err)
		} else {
			posts[i] = post
		}
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
