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

	val := trie.Get(key)
	if val != ih {
		t.Errorf("expected index, got %s", val)
	}

	key = "/a/deep/path"
	dh := testHandler{name: "deep"}
	_ = trie.Put(key, dh)

	val = trie.Get(key)
	if val != dh {
		t.Errorf("expected index, got %s", val)
	}

	val = trie.Get("/not/found")
	if val != nil {
		t.Errorf("expected no value, got %s", val)
	}
}

func TestGetStar(t *testing.T) {
	trie := NewTrie()

	key := "/static/*"
	h := testHandler{name: "static"}
	_ = trie.Put(key, h)

	key = "/static/img/logo.svg"
	val := trie.Get(key)
	if val != h {
		t.Errorf("expected 'static', got %v", val)
	}

	key = "/static/favicon.ico"
	val = trie.Get(key)
	if val != h {
		t.Errorf("expected 'static', got %v", val)
	}
}

func TestGetWithParams(t *testing.T) {
	trie := NewTrie()

	ch := testHandler{name: "commends"}

	key := "/accounts/:id/comments"
	_ = trie.Put(key, ch)

	val, params := trie.GetWithParams("/accounts/123/comments")
	if val != ch {
		t.Errorf("expected 'comments', got %s", val)
	}
	id, found := params["id"]
	if !found {
		t.Error("'id' parameter not found")
	}
	if id != "123" {
		t.Errorf("expected id=123, got %s", id)
	}

	ph := testHandler{name: "posts"}

	key = "/accounts/:id/posts/:filter"
	_ = trie.Put(key, ph)

	val, params = trie.GetWithParams("/accounts/314/posts/date")
	if val != ph {
		t.Errorf("expected 'posts', got %s", val)
	}
	id, found = params["id"]
	if !found {
		t.Error("'id' parameter not found")
	}
	if id != "314" {
		t.Errorf("expected id=314, got %s", id)
	}
	filter, found := params["filter"]
	if !found {
		t.Error("'filter' parameter not found")
	}
	if filter != "date" {
		t.Errorf("expected filter=date, got %s", filter)
	}

	sh := testHandler{name: "sign-in"}

	key = "/accounts/sign-in"
	_ = trie.Put(key, sh)

	val = trie.Get("/accounts/sign-in")
	if val != sh {
		t.Errorf("expected 'sign-in', got %s", val)
	}
}
