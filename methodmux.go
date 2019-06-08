package servemux

import (
	"net/http"
)

// MethodMux multiplexes HTTP requests by HTTP method.
type MethodMux map[string]http.Handler

func (m MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, found := m[r.Method]
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
