package servemux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func paramDumpHandler(w http.ResponseWriter, r *http.Request) {
	id := Value(r, "id")
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

	want := "123"
	got := w.Body.String()
	if got != want {
		t.Errorf("body = %q; want %q", got, want)
	}
}
