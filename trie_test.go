package servemux

import (
	"net/http"
	"testing"
)

type testHandler struct {
	name string
	http.Handler
}

var ih = testHandler{name: "index"}

func TestPut(t *testing.T) {
	trie := NewTrie()

	key := "/"
	isNew := trie.Put(key, ih)
	if !isNew {
		t.Errorf("expected key %s to be missing", key)
	}

	isNew = trie.Put(key, ih)
	if isNew {
		t.Errorf("expected key %s to have value", key)
	}
}

func TestGet(t *testing.T) {
	trie := NewTrie()

	key := "/"
	_ = trie.Put(key, ih)

	got, _ := trie.Get(key)
	if got != ih {
		t.Errorf("expected index, got %s", got)
	}

	key = "/a/deep/path"
	dh := testHandler{name: "deep"}
	_ = trie.Put(key, dh)

	got, _ = trie.Get(key)
	if got != dh {
		t.Errorf("expected index, got %s", got)
	}

	key = "/a/deep/path/"
	_ = trie.Put(key, dh)

	got, _ = trie.Get(key)
	if got != dh {
		t.Errorf("expected index, got %s", got)
	}

	got, _ = trie.Get("/not/found")
	if got != nil {
		t.Errorf("expected no value, got %s", got)
	}
}

func TestGetMatchAll(t *testing.T) {
	trie := NewTrie()

	key := "/static/*"
	h := testHandler{name: "static"}
	_ = trie.Put(key, h)

	key = "/static/img/logo.svg"
	got, args := trie.Get(key)
	if got != h {
		t.Errorf("expected 'static', got %v", got)
	}
	av := args["*"]
	if av != "img/logo.svg" {
		t.Errorf("expected 'img/logo.svg', got %v", av)
	}

	key = "/static/favicon.ico"
	got, args = trie.Get(key)
	if got != h {
		t.Errorf("expected 'static', got %v", got)
	}
	av = args["*"]
	if av != "favicon.ico" {
		t.Errorf("expected 'favicon.ico', got %v", av)
	}
}

func TestGetWithParams(t *testing.T) {
	trie := NewTrie()

	ch := testHandler{name: "commends"}

	key := "/accounts/:id/comments"
	_ = trie.Put(key, ch)

	got, args := trie.Get("/accounts/123/comments")
	if got != ch {
		t.Errorf("expected 'comments', got %s", got)
	}
	id, found := args["id"]
	if !found {
		t.Error("'id' parameter not found")
	}
	if id != "123" {
		t.Errorf("expected id=123, got %s", id)
	}

	ph := testHandler{name: "posts"}

	key = "/accounts/:id/posts/:filter"
	_ = trie.Put(key, ph)

	got, args = trie.Get("/accounts/314/posts/date")
	if got != ph {
		t.Errorf("expected 'posts', got %s", got)
	}
	id, found = args["id"]
	if !found {
		t.Error("'id' parameter not found")
	}
	if id != "314" {
		t.Errorf("expected id=314, got %s", id)
	}
	filter, found := args["filter"]
	if !found {
		t.Error("'filter' parameter not found")
	}
	if filter != "date" {
		t.Errorf("expected filter=date, got %s", filter)
	}

	sh := testHandler{name: "sign-in"}

	key = "/accounts/sign-in"
	_ = trie.Put(key, sh)

	got, _ = trie.Get("/accounts/sign-in")
	if got != sh {
		t.Errorf("expected 'sign-in', got %s", got)
	}
}
