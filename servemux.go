package servemux

import (
	"net/http"
	"sync"
)

// ServeMux to be defined
type ServeMux struct {
	mu   sync.RWMutex
	trie *Trie
}

// New allocates and returns a new ServeMux.
func New() *ServeMux {
	return &ServeMux{
		trie: NewTrie(),
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

	if h == nil {
		return notFoundHandler, "" // TODO: something better needed.
	}

	return h, ""
}

// NotFoundHandler to be defined.
type NotFoundHandler struct {
}

// ServeHTML to be defined.
func (h NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: something better!
	w.Write([]byte("Not found"))
}

var notFoundHandler = NotFoundHandler{}
