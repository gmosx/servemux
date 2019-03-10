package servemux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getHandler(t *testing.T, mux *ServeMux, path string) http.Handler {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	h, _ := mux.Handler(req)
	return h
}

func TestMux(t *testing.T) {
	ih := testHandler{name: "index"}
	dh := testHandler{name: "deep"}
	ph := testHandler{name: "posts"}
	ch := testHandler{name: "comments"}

	mux := New()
	mux.Handle("/", ih)
	mux.Handle("/a/deep/path", dh)
	mux.Handle("/accounts/:id/posts", ph)
	mux.Handle("/accounts/:id/comments", ch)

	h := getHandler(t, mux, "/")
	if h != ih {
		t.Errorf("Expected 'index', got %s", h)
	}

	h = getHandler(t, mux, "/a/deep/path")
	if h != dh {
		t.Errorf("Expected 'deep', got %s", h)
	}

	h = getHandler(t, mux, "/invalid/path")
	if h != notFoundHandler {
		t.Errorf("Expected notFoundHandler, got %s", h)
	}

	h = getHandler(t, mux, "/accounts/123/comments")
	if h != ch {
		t.Errorf("Expected 'comments', got %s", h)
	}

	h = getHandler(t, mux, "/accounts/123/posts")
	if h != ph {
		t.Errorf("Expected 'comments', got %s", h)
	}
}

func paramDumpHandler(w http.ResponseWriter, r *http.Request) {
	id, found := ParamValue(r, "id")
	if found {
		fmt.Fprintf(w, id)
	}
}
func TestServeHTTP(t *testing.T) {
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
