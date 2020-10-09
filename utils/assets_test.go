package utils

import (
	"testing"
)

// Ensure that assetsRoot points to a real filepath.
func TestAssetsRootExists(t *testing.T) {
	if ar, err := assetsRoot(); err != nil {
		t.Fatalf("calling assetsRoot(): %+v", err)
	} else if ar == "" {
		t.Fatal("Assets root not found")
	}
	// Check that it caches correctly
	if cachedAssetsRoot == "" {
		t.Fatal("Assets root did not cache")
	}
}
