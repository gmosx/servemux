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

	val, _ := trie.Get(key)
	if val != ih {
		t.Errorf("expected index, got %s", val)
	}

	key = "/a/deep/path"
	dh := testHandler{name: "deep"}
	_ = trie.Put(key, dh)

	val, found := trie.Get(key)
	if !found {
		t.Errorf("expected to find %s", key)
	}
	if val != dh {
		t.Errorf("expected index, got %s", val)
	}

	val, found = trie.Get("/not/found")
	if found {
		t.Errorf("expected no value, got %s", val)
	}
}
