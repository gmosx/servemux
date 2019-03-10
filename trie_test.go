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

func TestGetWithParams(t *testing.T) {
	trie := NewTrie()

	key := "/accounts/:id/comments"
	_ = trie.Put(key, ih)

	val, params := trie.GetWithParams("/accounts/123/comments")
	if val != ih {
		t.Errorf("expected 'index', got %s", val)
	}
	id, found := params["id"]
	if !found {
		t.Error("'id' parameter not found")
	}
	if id != "123" {
		t.Errorf("expected id=123, got %s", id)
	}
}
