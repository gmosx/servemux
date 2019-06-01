package servemux

import (
	"net/http"
)

// MethodMux multiplexes HTTP requests by HTTP method.
type MethodMux struct {
	handlers map[string]http.Handler
}

// MuxMethods makes a MethodMux from a variadic arguments list.
func MuxMethods(args ...interface{}) *MethodMux {
	handlers := map[string]http.Handler {}

	var meth string

	for i, p := range args {
		if i % 2 == 0 {
			m, ok := p.(string)
			if !ok {
				panic("invalid arguments to MuxMethods")
			}
			meth = m
			continue
		}
		switch h := p.(type) {
		case http.Handler:
			handlers[meth] = h
		case func(http.ResponseWriter, *http.Request):
			handlers[meth] = http.HandlerFunc(h)
		}
	} 

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
