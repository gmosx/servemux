package servemux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodMuxAllow(t *testing.T) {
	geth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET"))
	})

	deleteh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("DELETE"))
	})

	mux := MethodMux{
		http.MethodGet: geth,
		http.MethodDelete: deleteh,
	}

	allow := []string{"GET", "DELETE"}

	for _, m := range allow {
		w := httptest.NewRecorder()

		r, err := http.NewRequest(m, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		mux.ServeHTTP(w, r)

		want := m
		got := w.Body.String()
		if got != want {
			t.Errorf("body = %q; want %q", got, want)
		}
	}
}

func TestMethodMuxReject(t *testing.T) {
	geth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET"))
	})

	mux := MethodMux{
		http.MethodGet: geth,
	}

	w := httptest.NewRecorder()

	r, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	mux.ServeHTTP(w, r)

	want := 404
	got := w.Code
	if got != want {
		t.Errorf("body = %d; want %d", got, want)
	}
}
