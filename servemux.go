package servemux

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type key struct{}

var paramsKey key

// ServeMux is an HTTP request multiplexer. It matches the URL of each incoming
// request against a list of registered patterns and calls the handler for the
// pattern that most closely matches the URL.
//
// ServeMux is a minimal extension of http.ServeMux found in the standard
// library. It offers improved performance and more powerful pattern-marching
// (e.g. parameters and match-all options).
type ServeMux struct {
	mu              sync.RWMutex
	trie            *Trie
	NotFoundHandler http.Handler
}

// New allocates and returns a new ServeMux.
func New() *ServeMux {
	return &ServeMux{
		trie:            NewTrie(),
		NotFoundHandler: http.HandlerFunc(http.NotFound),
	}
}

// Handle registers the handler for the given pattern. If a handler already
// exists for pattern, Handle panics.
func (m *ServeMux) Handle(pattern string, handler http.Handler) {
	m.mu.Lock()
	newval := m.trie.Put(pattern, handler)
	if !newval {
		panic(fmt.Sprintf("Duplicate handler for pattern '%s'", pattern))
	}
	m.mu.Unlock()
}

// HandleFunc registers the handler function for the given pattern.
func (m *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}

	m.Handle(pattern, http.HandlerFunc(handler))
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	h, params := m.trie.Get(p)

	if h == nil {
		m.NotFoundHandler.ServeHTTP(w, r)
		return
	}

	if params != nil {
		ctx := context.WithValue(r.Context(), paramsKey, params)
		h.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	h.ServeHTTP(w, r)
}

// ParamValue returns the value associated with key.
func ParamValue(r *http.Request, key string) string {
	params := r.Context().Value(paramsKey)
	if params == nil {
		return ""
	}

	v, found := params.(map[string]string)[key]
	if !found {
		return ""
	}

	return v
}
