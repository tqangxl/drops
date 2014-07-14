/*
Package tst is a pure Go implementation of *Ternary Search Trie*.

It's also called *Ternary Search Tree*, here's a detailed
description by the inventors Jon Bently and Robert Sedgwick,
http://www.drdobbs.com/database/ternary-search-trees/184410528.

It can be used as an efficient symbol table (as efficient as hash table), and
supports order traversal and many advanced search operations like prefix and
wildcard searches.
*/
package element

// import "fmt"

type node struct {
	c     byte
	left  *node
	right *node
	mid   *node
	val   *View
}

// Trie represents a Ternary Search Trie. An empty value can be used to insert
// key/value pairs.
type ViewTrie struct {
	root *node
	n    int
}

// Size returns number of elements in the trie.
func (t *ViewTrie) Size() int {
	return t.n
}

// Put inserts a key/value pair into the trie. If the key is already in the
// trie, it's value will be updated. Note: empty strings will be ignored.
func (t *ViewTrie) Put(key string, value *View) {
	t.put(key, value, false)
}

// PutRK is the same as Put, but the key is reversely inserted into the
// trie.
func (t *ViewTrie) PutRK(key string, value *View) {
	t.put(key, value, true)
}

// Get returns the value associated with key.
func (t *ViewTrie) Get(key string) *View {
	return t.get(key, false)
}

// GetRK is the same as Get, but the key is reversely iterated when searching
// for the value.
func (t *ViewTrie) GetRK(key string) *View {
	return t.get(key, true)
}

// GetShortestPrefix searches for the shortest key which is a prefix of the
// given key and returns the value. For example, if trie contains key "com"
// and "com.g", search for "com.google" will return value associated with
// "com".
func (t *ViewTrie) GetShortestPrefix(key string) *View {
	return t.getShortestPrefix(key, false)
}

// GetShortestPrefixRK is the same as GetShortestPrefix, but the key is
// reversely iterated when searching for the value.
func (t *ViewTrie) GetShortestPrefixRK(key string) *View {
	return t.getShortestPrefix(key, true)
}

func iterStartEndStep(key string, reverseKey bool) (start, end, step int) {
	if reverseKey {
		return len(key) - 1, 0, -1
	} else {
		return 0, len(key) - 1, 1
	}
}

func (t *ViewTrie) put(key string, value *View, reverseKey bool) {
	if len(key) < 1 {
		return
	}

	start, end, step := iterStartEndStep(key, reverseKey)

	t.n++
	pnd := &t.root
	i := start
	c := key[start]
	for {
		if *pnd == nil {
			// fmt.Printf("new node: %c\n", c)
			*pnd = &node{c: c}
		}
		switch {
		case c < (*pnd).c:
			pnd = &(*pnd).left
		case c > (*pnd).c:
			pnd = &(*pnd).right
		case i != end:
			pnd = &(*pnd).mid
			i += step
			c = key[i]
		default:
			(*pnd).val = value
			// fmt.Println("add value:", value)
			return
		}
	}
}

func (t *ViewTrie) get(key string, reverseKey bool) *View {
	if len(key) < 1 {
		return nil
	}

	start, end, step := iterStartEndStep(key, reverseKey)

	nd := t.root
	i := start
	c := key[start]
	// go down the tree
	for nd != nil {
		switch {
		case c < nd.c:
			nd = nd.left
		case c > nd.c:
			nd = nd.right
		case i != end:
			nd = nd.mid
			i += step
			c = key[i]
		default:
			return nd.val
		}
	}
	return nil
}

func (t *ViewTrie) getShortestPrefix(key string, reverseKey bool) *View {
	// Most code copied from Get.
	if len(key) < 1 {
		return nil
	}

	start, end, step := iterStartEndStep(key, reverseKey)

	nd := t.root
	i := start
	c := key[start]
	for nd != nil {
		switch {
		case c < nd.c:
			nd = nd.left
		case c > nd.c:
			nd = nd.right
		// This is the only added code to get.
		case nd.val != nil:
			return nd.val
		case i != end:
			nd = nd.mid
			i += step
			c = key[i]
		default:
			return nd.val
		}
	}
	return nil
}
