# ServeMux

```go
import "github.com/gmosx/go-servemux"

mux := servemux.New()
mux.Handle("/about", NewAboutHandler())

log.Fatal(http.ListenAndServe(":8080", mux))
```

## License

MIT License.