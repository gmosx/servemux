package servemux

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type key struct{}

var argsKey key

// ServeMux is an HTTP request multiplexer. It matches the URL of each incoming
// request against a list of registered patterns and calls the handler for the
// pattern that most closely matches the URL.
//
// ServeMux is a minimal extension of http.ServeMux found in the standard
// library. It offers improved performance and parameterized pattern-marching.
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
	h, args := m.trie.Get(p)

	if h == nil {
		m.NotFoundHandler.ServeHTTP(w, r)
		return
	}

	if args != nil {
		ctx := context.WithValue(r.Context(), argsKey, args)
		h.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	h.ServeHTTP(w, r)
}

// Value returns the argument value associated with key.
func Value(r *http.Request, key string) string {
	args := r.Context().Value(argsKey)
	if args == nil {
		return ""
	}

	v, found := args.(map[string]string)[key]
	if !found {
		return ""
	}

	return v
}
