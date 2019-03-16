package servemux

import (
	"net/http"
	"strings"
)

// Value represents a trie value.
type Value = http.Handler

// Trie is a prefix search tree, specialized to work with the ServeMux.
type Trie struct {
	value    Value
	special  string
	children map[string]*Trie // TODO: strange to use map within a trie :-|
}

// NewTrie allocates and returns a new *Trie.
func NewTrie() *Trie {
	return &Trie{
		children: make(map[string]*Trie),
	}
}

// Put inserts a new value into the tree. Returns true if it inserts a new value,
// false if it replaces an existing value.
func (t *Trie) Put(key string, val Value) bool {
	node := t
	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		if len(part) != 0 { // TODO: remove this test?
			// Check if the part is 'special' e.g. a parameter or a match-all
			// pattern.
			if part[0] == '*' {
				node.special = part
				break
			}

			if part[0] == ':' {
				node.special = part
			}
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

// Get returns the value associated with the given key.
func (t *Trie) Get(key string) Value {
	node := t
	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		node, _ = selectChild(node, part)
		if node == nil {
			return nil
		}
		if node.special == "*" {
			break
		}
		if i == -1 {
			break
		}
	}

	return node.value
}

// GetWithParams returns the value associated with the given key.
func (t *Trie) GetWithParams(key string) (Value, map[string]string) {
	var params map[string]string
	node := t
	for part, i := splitter(key, 0); ; part, i = splitter(key, i) {
		n, isParamMatch := selectChild(node, part)
		if n == nil {
			return nil, params
		}
		if isParamMatch {
			if params == nil {
				params = map[string]string{
					node.special[1:]: part,
				}
			} else {
				params[node.special[1:]] = part
			}
		}
		node = n
		if node.special == "*" {
			break
		}
		if i == -1 {
			break
		}
	}

	return node.value, params
}

func splitter(path string, start int) (segment string, next int) {
	if /* len(path) == 0 || start < 0 || */ start > len(path)-1 {
		return "", -1
	}

	end := strings.IndexRune(path[start+1:], '/')
	if end == -1 {
		return path[start+1:], -1
	}

	return path[start+1 : start+end+1], start + end + 1
}

func selectChild(node *Trie, key string) (*Trie, bool) {
	c, found := node.children[key]
	if found {
		return c, false
	}

	if node.special != "" {
		return node.children[node.special], true
	}

	return nil, false
}
