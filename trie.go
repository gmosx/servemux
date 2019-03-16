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
	pattern  string
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
	for segment, i := sliceSegmentAt(key, 0); ; segment, i = sliceSegmentAt(key, i) {
		if len(segment) != 0 { // TODO: remove this test?
			// Check if the segment is a parameter or a match-all pattern.
			if segment[0] == '*' {
				node.pattern = segment
				break
			}

			if segment[0] == ':' {
				node.pattern = segment
			}
		}

		child, _ := node.children[segment]
		if child == nil {
			child = NewTrie()
			node.children[segment] = child
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
	for segment, i := sliceSegmentAt(key, 0); ; segment, i = sliceSegmentAt(key, i) {
		node, _ = selectChild(node, segment)
		if node == nil {
			return nil
		}
		if node.pattern == "*" {
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
	for segment, i := sliceSegmentAt(key, 0); ; segment, i = sliceSegmentAt(key, i) {
		n, isParamMatch := selectChild(node, segment)
		if n == nil {
			return nil, params
		}
		if isParamMatch {
			if params == nil {
				params = map[string]string{
					node.pattern[1:]: segment,
				}
			} else {
				params[node.pattern[1:]] = segment
			}
		}
		node = n
		if node.pattern == "*" {
			break
		}
		if i == -1 {
			break
		}
	}

	return node.value, params
}

// slice returns the path segment that begins at start and the index
// for the next segment (or -1 if the end is reached).
func sliceSegmentAt(path string, start int) (segment string, next int) {
	if /* len(path) == 0 || start < 0 || */ start > len(path)-1 {
		return "", -1
	}

	end := strings.IndexRune(path[start+1:], '/')
	if end == -1 {
		return path[start+1:], -1
	}

	return path[start+1 : start+end+1], start + end + 1
}

// selectChild selects the next child for inspection. A second boolean return-value
// is true if a pattern is used for selection.
func selectChild(node *Trie, key string) (*Trie, bool) {
	c, found := node.children[key]
	if found {
		return c, false
	}

	if node.pattern != "" {
		return node.children[node.pattern], true
	}

	return nil, false
}
