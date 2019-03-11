package servemux

import (
	"context"
	"net/http"
	"sync"
)

type key struct{}

var paramsKey key

// ServeMux to be defined
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

// Handle to be defined
func (m *ServeMux) Handle(pattern string, handler http.Handler) {
	m.mu.Lock()
	m.trie.Put(pattern, handler)
	m.mu.Unlock()
}

// HandleFunc to be defined
func (m *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}

	m.Handle(pattern, http.HandlerFunc(handler))
}

// Handler returns the handler to use for the given request, consulting r.Method,
// r.Host, and r.URL.Path. It always returns a non-nil handler.
func (m *ServeMux) Handler(r *http.Request) (handler http.Handler, pattern string) {
	p := r.URL.Path
	h := m.trie.Get(p)

	// TODO: return matched pattern string.

	if h == nil {
		return m.NotFoundHandler, ""
	}

	return h, ""
}

func (m *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	h, params := m.trie.GetWithParams(p)

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
