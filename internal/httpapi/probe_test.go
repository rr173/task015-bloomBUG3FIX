package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"task015-bloom/internal/bloom"
)

// TestProbeAddNullItem verifies that POST /add with a null or missing "item"
// field returns 400 instead of panicking with a nil pointer dereference.
func TestProbeAddNullItem(t *testing.T) {
	f, _ := bloom.New(1000, 0.01)
	srv := httptest.NewServer(New(f).Handler())
	defer srv.Close()

	cases := []struct {
		name string
		body string
	}{
		{"null item", `{"item":null}`},
		{"missing item", `{}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, err := http.Post(srv.URL+"/add", "application/json", strings.NewReader(c.body))
			if err != nil {
				t.Fatalf("request failed (server likely panicked): %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				t.Errorf("status = 200, want non-200 error response for nil item")
			}
		})
	}
}
