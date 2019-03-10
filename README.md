# Nitro ServeMux

import "github.com/gmosx/go-servemux"


```go
mux := servemux.New()
mux.Handle("/about", NewAboutHandler())
mux.Handle("/blog", NewStaticHandler())

log.Fatal(http.ListenAndServe(":8080", mux))
```