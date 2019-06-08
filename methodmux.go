package servemux

import (
	"net/http"
	"strings"
)

// MethodMux multiplexes HTTP requests by HTTP method.
type MethodMux map[string]http.Handler

func (m MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, found := m[r.Method]
	if !found {
		m.methodNotAllowed(w, r)
		// http.NotFound(w, r)
		return
	}

	h.ServeHTTP(w, r)
}

func (m MethodMux) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if len(m) > 0 {
		allowMethods := make([]string, len(m))

		i := 0
		for k := range m {
			allowMethods[i] = k
			i++
		}

		w.Header().Set("Allow", strings.Join(allowMethods, ", "))
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
