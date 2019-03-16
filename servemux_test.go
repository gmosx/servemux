package servemux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func paramDumpHandler(w http.ResponseWriter, r *http.Request) {
	id := ParamValue(r, "id")
	fmt.Fprintf(w, id)
}

func TestServeHTTPWithParams(t *testing.T) {
	mux := New()
	mux.HandleFunc("/accounts/:id/posts", paramDumpHandler)

	r, err := http.NewRequest("GET", "/accounts/123/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	s := w.Body.String()
	if s != "123" {
		t.Errorf("Expected 123, got %s", s)
	}
}
