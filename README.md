# ServeMux

`ServeMux` is a drop-in replacement of the http.ServeMux in the standard library. It offers additional functionality while staying API compatible.

## Example

```go
import "github.com/gmosx/go-servemux"

func postsHandler(w http.ResponseWriter, r *http.Request) {
    id := ParamValue(r, "id")
    fmt.Fprintf(w, id)
}

mux := servemux.New()
mux.HandleFunc("/accounts/:id/posts", postsHandler)

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