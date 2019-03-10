package servemux

import (
	"net/http"
	"strings"
)

// Value represents a trie value.
type Value = http.Handler

var nilValue = NotFoundHandler{}

// Trie is a prefix search tree.
type Trie struct {
	value    Value
	param    string
	children map[string]*Trie // TODO: strange to use map within a trie :-|
}

func splitter(path string, start int) (segment string, next int) {
	// TODO: remove some overzealous checks?
	if len(path) == 0 || start < 0 || start > len(path)-1 {
		return "", -1
	}

	end := strings.IndexRune(path[start+1:], '/')
	if end == -1 {
		return path[start+1:], -1
	}

	return path[start+1 : start+end+1], start + end + 1
}

// NewTrie allocates and returns a new *Trie.
func NewTrie() *Trie {
	return &Trie{
		children: make(map[string]*Trie),
	}
}

func isParam(key string) bool {
	return strings.HasPrefix(key, ":") // TODO: optimize me!
}

// Put inserts a new value into the tree.
func (t *Trie) Put(key string, val Value) bool {
	node := t

	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		if isParam(part) {
			node.param = part
		}

		child, _ := node.children[part]
		if child == nil {
			child = NewTrie()
			node.children[part] = child
		}
		node = child
		if i == -1 {
			break
		}
	}

	isNewVal := node.value == nil
	node.value = val
	return isNewVal
}

func selectChild(node *Trie, key string) *Trie {
	c, found := node.children[key]
	if found {
		return c
	}

	if node.param != "" {
		return node.children[node.param]
	}

	return nil
}

// Get returns the value associated with the given key.
func (t *Trie) Get(key string) (Value, bool) {
	node := t
	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		node = selectChild(node, part)
		if node == nil {
			return nilValue, false
		}
		if i == -1 {
			break
		}
	}

	return node.value, true
}
