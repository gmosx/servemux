package main

import (
	"log"
	"net/http"

	"go.reizu.org/servemux"
)

func main() {
	mux := servemux.New()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		w.Write([]byte("Welcome!\n"))
	})

	mux.HandleFunc("/user/:id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		w.Write([]byte(servemux.Value(r, "id")))
	})

	mux.Handle("/post/:id", servemux.NewMethodMux(map[string]http.Handler{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("GET!\n"))
		}),
		http.MethodDelete: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("DELETE!\n"))
		}),
	}))

	log.Fatal(http.ListenAndServe(":3000", mux))
}
