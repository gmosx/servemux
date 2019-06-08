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

	getFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("GET!\n"))
	}

	deleteFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("DELETE!\n"))
	}

	mux.Handle("/post/:id", servemux.MethodMux{
		http.MethodGet:    http.HandlerFunc(getFunc),
		http.MethodDelete: http.HandlerFunc(deleteFunc),
	})

	log.Fatal(http.ListenAndServe(":3000", mux))
}
