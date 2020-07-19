package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

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
