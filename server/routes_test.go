package server

import (
	"fmt"
	"testing"
	"net/http"
	"net/http/httptest"
)

// Ensures that certain "standalone" pages serve 200 responses.
func TestStandalonePages(t *testing.T) {
	srv, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	srv.SetupRoutes()
	base := fmt.Sprintf("http://localhost:%d", srv.Port)
	for _, route := range []string{"/", "/about", "/now"} {
		url := fmt.Sprintf("%s%s", base, route)
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, r)
		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("url %s: got status code %d; want %d", url, w.Result().StatusCode, http.StatusOK)
		}
	}
}
