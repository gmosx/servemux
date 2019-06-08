package servemux

import (
	"net/http"
)

// MethodHandlers multiplexes HTTP requests by HTTP method.
type MethodHandlers map[string]http.Handler

func (m MethodHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, found := m[r.Method]
	if !found {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.ServeHTTP(w, r)
}

// MethodFuncs multiplexes HTTP requests by HTTP method.
type MethodFuncs map[string]func(http.ResponseWriter, *http.Request)

func (m MethodFuncs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, found := m[r.Method]
	if !found {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h(w, r)
}

// func (m MethodHandlers) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
// 	if len(m) > 0 {
// 		allowMethods := make([]string, len(m))

// 		i := 0
// 		for k := range m {
// 			allowMethods[i] = k
// 			i++
// 		}

// 		w.Header().Set("Allow", strings.Join(allowMethods, ", "))
// 	}
// 	w.WriteHeader(http.StatusMethodNotAllowed)
// }
