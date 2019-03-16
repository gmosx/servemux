package main

import (
	"log"
	"net/http"

	"go.reizu.org/pkg/servemux"
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

	log.Fatal(http.ListenAndServe(":3000", mux))
}
