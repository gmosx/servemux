package servemux

import (
	"net/http"
	"strings"
)

// TrieValue represents a trie value.
type TrieValue = http.Handler

// Trie is a prefix search tree, specialized to work with ServeMux.
type Trie struct {
	value    TrieValue
	param    string
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
func (t *Trie) Put(key string, val TrieValue) bool {
	node := t

	for segment, i := sliceSegmentAt(key, 0); ; segment, i = sliceSegmentAt(key, i) {
		if len(segment) != 0 { // TODO: remove this test?
			// Check if the segment is a parameter.
			if segment[0] == '*' {
				node.param = segment
				break
			}

			if segment[0] == ':' {
				node.param = segment
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

// Get returns the value associated with the given key and optionally
// a map of arguments.
func (t *Trie) Get(key string) (TrieValue, map[string]string) {
	var args map[string]string

	node := t
	for segment, i := sliceSegmentAt(key, 0); ; segment, i = sliceSegmentAt(key, i) {
		child := selectChild(node, segment)
		if child == nil {
			return nil, args
		}

		if node.param != "" {
			if args == nil {
				args = map[string]string{}
			}
			if node.param[0] == ':' {
				args[node.param[1:]] = segment
			}
		}

		node = child

		if node.param == "*" {
			if args == nil {
				args = map[string]string{}
			}
			args["*"] = key[i+1:]
			break
		}

		if i == -1 {
			break
		}
	}

	return node.value, args
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
