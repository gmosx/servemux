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
	children map[string]*Trie
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

	for segment, i := sliceSegmentAt(key, 1); ; segment, i = sliceSegmentAt(key, i) {
		if len(segment) != 0 {
			if segment[0] == '*' {
				node.param = segment
				i = -1
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

	newval := node.value == nil
	node.value = val
	return newval
}

// Get returns the value associated with the given key and optionally
// a map of arguments.
func (t *Trie) Get(key string) (TrieValue, map[string]string) {
	var args map[string]string

	node := t
	prev := 1
	for segment, i := sliceSegmentAt(key, prev); ; segment, i = sliceSegmentAt(key, i) {
		child := selectChild(node, segment)

		if child == nil {
			return nil, args
		}

		if node.param != "" {
			if args == nil {
				args = map[string]string{}
			}
			if node.param[0] == '*' {
				args["*"] = key[prev:]
				i = -1
			}
			if node.param[0] == ':' {
				args[node.param[1:]] = segment
			}
		}

		node = child

		if i == -1 {
			break
		}

		prev = i
	}

	return node.value, args
}

// slice returns the path segment that begins at start and the index
// for the next segment (or -1 if the end is reached).
func sliceSegmentAt(path string, start int) (segment string, next int) {
	if /* len(path) == 0 || start < 0 || */ start > len(path)-1 {
		return "", -1
	}

	end := strings.IndexRune(path[start:], '/')
	if end == -1 {
		return path[start:], -1
	}

	return path[start : start+end], start + end + 1
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
