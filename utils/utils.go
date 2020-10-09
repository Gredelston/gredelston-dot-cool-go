package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Panicf formats s, and panics.
func Panicf(s string, i ...interface{}) { panic(fmt.Sprintf(s, i...)) }

// PathExists checks whether a fully qualified path is present on the server's filesystem.
func PathExists(fp string) (bool, error) {
	_, err := os.Stat(fp)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var cachedStaticRoot string

// staticRoot finds the directory containing project's static files.
func staticRoot() string {
	if cachedStaticRoot != "" {
		return cachedStaticRoot
	}
	var attempted []string
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		d := filepath.Join(gopath, "src", "github.com", GitUserName, ProjectName, "static")
		if _, err := os.Stat(d); !os.IsNotExist(err) {
			cachedStaticRoot = d
			return cachedStaticRoot
		}
		attempted = append(attempted, d)
	}
	Panicf("static dir not found for gopath=%q; attempted %+v", os.Getenv("GOPATH"), attempted)
	return ""
}

// StaticPath finds the path to a specified static file or directory.
func StaticPath(elem ...string) string {
	return filepath.Join(append([]string{staticRoot()}, elem...)...)
}
