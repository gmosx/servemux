# ServeMux

`ServeMux` is an efficient, more flexible, yet API-compatible extension of `http.ServeMux`.

## Features

* Uses a specialized [trie](https://en.wikipedia.org/wiki/Trie) data-structure for efficiency
* Parameterized pattern matching (:segment, *)
* Orthogonal multiplexing by request method
* Simple, drop-in replacement for `http.ServeMux`
* No external dependencies
* No extraneous features

## Example

```go
import "go.reizu.org/servemux"

func postsHandler(w http.ResponseWriter, r *http.Request) {
    id := servemux.Value(r, "id")
    fmt.Fprintf(w, id)
}

mux := servemux.New()

// Example matches:
// /accounts/1/posts
// /accounts/2/posts
mux.HandleFunc("/accounts/:id/posts", postsHandler)

// Example matches:
// /static/img/logo.png
// /static/favicon.ico
// The matched sub-path can be accessed with: servemux.Value(r, "*")
mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

// Multiplex multiple handlers by the request method:
mux.Handle("/post/:id", servemux.ByMethod(
    http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("GET!\n"))
    },
    http.MethodDelete, func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("DELETE!\n"))
    },
))

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

## FAQ

### Why doesn't ServeMux match for the request method?

`ServeMux` routes URLs to handlers. In REST terms, it serves Web Resources at specific URLs. The request method is an orthogonal concern best handled within the handler itself. This way, the conceptual simplicity of Go (along with compatibility with the standard library) is retained and you can easily reuse code for different methods within the handler.

### That's reasonable, but I still want to multiplex by request method

We provide a `MethodMux` handler to do just that. For convenience, you can use the `servemux.ByMethod` helper. Check out the included example for more details.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please make sure to update tests as appropriate.

## Contact

[@gmosx](https://twitter.com/gmosx) on Twitter.

## License

MIT, see [LICENSE](./LICENSE) file for details.

Copyright 2019 George Moschovitis.
