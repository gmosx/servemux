# ServeMux

`ServeMux` is a drop-in replacement of the http.ServeMux in the standard library. It offers additional functionality while staying API compatible.

## Example

```go
import "go.reizu.org/pkg/servemux"

func postsHandler(w http.ResponseWriter, r *http.Request) {
    id := servemux.Value(r, "id")
    fmt.Fprintf(w, id)
}

mux := servemux.New()
mux.HandleFunc("/accounts/:id/posts", postsHandler)
mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

log.Fatal(http.ListenAndServe(":8080", mux))
```

## Benchmarking

Run the example:

```sh
go run _example/main.go
```

Use [bombardier](https://github.com/codesenberg/bombardier) to benchmark the performance of servemux:

```sh
./bombardier -c 125 -n 1000000 http://localhost:3000/
./bombardier -c 125 -n 1000000 http://localhost:3000/user/23
```

## Contact

[@gmosx](https://twitter.com/gmosx) on Twitter.

## License

MIT, see `LICENSE` file for details.