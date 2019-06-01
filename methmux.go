package servemux

import (
	"net/http"
)

// MethodMux multiplexes HTTP requests by HTTP method.
type MethodMux struct {
	handlers map[string]http.Handler
}

// NewMethodMux returns a pointer to a new MethodMux.
func NewMethodMux(handlers map[string]http.Handler) *MethodMux {
	return &MethodMux{handlers: handlers}
}
func (m *MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, found := m.handlers[r.Method]
	if !found {
		// allowMethods := make([]string, len(m.handlers))
		// i := 0
		// for k := range m.handlers {
		// 	allowMethods[i] = k
		// 	i++
		// }

		// if len(allowMethods) != 0 {
		// 	w.Header().Set("Allow", strings.Join(allowMethods, ", "))
		// }
		// w.WriteHeader(http.StatusMethodNotAllowed)
		http.NotFound(w, r)
		return
	}

	h.ServeHTTP(w, r)
}
