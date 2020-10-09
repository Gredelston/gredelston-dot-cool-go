package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var cachedAssetsRoot string

// assetsRoot finds the directory containing project's assets.
func assetsRoot() (string, error) {
	if cachedAssetsRoot != "" {
		return cachedAssetsRoot, nil
	}
	var attempted []string
	gopath := os.Getenv("GOPATH")
	for _, p := range strings.Split(gopath, ":") {
		ar := filepath.Join(p, "src", "github.com", GitUserName, ProjectName, "assets")
		if exists, err := PathExists(ar); err != nil {
			return "", err
		} else if exists {
			cachedAssetsRoot = ar
			return cachedAssetsRoot, nil
		}
		attempted = append(attempted, ar)
	}
	// If nothing in gopath worked, try whatever's local.
	ar, _ := filepath.Abs("assets")
	if exists, _ := PathExists(ar); exists {
		cachedAssetsRoot = ar
		return cachedAssetsRoot, nil
	}
	attempted = append(attempted, ar)
	return "", fmt.Errorf("assets dir not found with gopath=%q, attempted=%+v", gopath, attempted)
}

// Asset finds the path to a specified assets file or directory.
func Asset(elem ...string) (string, error) {
	ar, err := assetsRoot()
	if err != nil {
		return "", errors.New("determining assets root")
	}
	return filepath.Join(append([]string{ar}, elem...)...), nil
}
