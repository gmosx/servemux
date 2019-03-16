# ServeMux

`ServeMux` is an efficient, API-compatible extension of `http.ServeMux`.

## Features

* Uses a specialized [trie](https://en.wikipedia.org/wiki/Trie) data-structure for efficiency.
* Parameterized pattern matching (:segment, *)
* Simple, drop-in replacement for `http.ServeMux`
* No external dependencies

## Example

```go
import "go.reizu.org/pkg/servemux"

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

MIT, see [LICENSE](./LICENSE) file for details.

Copyright 2019 George Moschovitis.
