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

var cachedAssetsRoot string

// assetsRoot finds the directory containing project's assets.
func assetsRoot() string {
	if cachedAssetsRoot != "" {
		return cachedAssetsRoot
	}
	var attempted []string
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), ":") {
		d := filepath.Join(gopath, "src", "github.com", GitUserName, ProjectName, "assets")
		if _, err := os.Stat(d); !os.IsNotExist(err) {
			cachedAssetsRoot = d
			return cachedAssetsRoot
		}
		attempted = append(attempted, d)
	}
	Panicf("Assets dir not found for gopath=%q; attempted %+v", os.Getenv("GOPATH"), attempted)
	return ""
}

// Asset finds the path to a specified assets file or directory.
func Asset(elem ...string) string {
	return filepath.Join(append([]string{assetsRoot()}, elem...)...)
}
