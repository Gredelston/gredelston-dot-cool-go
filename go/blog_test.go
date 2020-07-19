package main

import "testing"

func TestBlogRootExists(t *testing.T) {
	s, err := NewServer()
	if err != nil { t.Error(err) }
	blogRoot := s.BlogRoot()
	if exists, err := s.FullPathExists(blogRoot); err != nil {
		t.Error(err)
	} else if !exists {
		t.Errorf("Blog root %q not present on server", blogRoot)
	}
}
