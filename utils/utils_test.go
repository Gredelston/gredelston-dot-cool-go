package utils

import (
	"testing"
)

// Ensures that staticRoot points to a real filepath.
func TestStaticRootExists(t *testing.T) {
	if staticRoot() == "" {
		t.Fatal("Static root not found")
	}
	// Check that it caches correctly
	if cachedStaticRoot == "" {
		t.Fatal("Static root did not cache")
	}
}
