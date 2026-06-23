package vulcain

import (
	"strings"

	"github.com/dunglas/httpsfv"
)

// node represents a node of a JSON document
type node struct {
	preload       bool
	preloadParams []*httpsfv.Params
	fields        bool
	fieldsParams  []*httpsfv.Params
	path          string
	parent        *node
	children      []*node
}

// _type is the type of operation to apply, can be Preload or Fields
type _type int

const (
	preload _type = iota
	fields
)

// importPointers imports JSON pointers in the tree
func (n *node) importPointers(t _type, pointers httpsfv.List) {
	for _, member := range pointers {
		// Ignore invalid value
		member, ok := member.(httpsfv.Item)
		if !ok {
			continue
		}

		pointer, ok := member.Value.(string)
		if !ok {
			continue
		}

		pointer = strings.Trim(pointer, "/")
		if pointer != "" {
			partsToTree(t, strings.Split(pointer, "/"), n, member.Params)
		}
	}
}

// String returns a JSON pointer
func (n *node) String() string {
	if n.parent == nil {
		return "/"
	}

	s := n.path
	c := n.parent
	for c != nil {
		s = c.path + "/" + s
		c = c.parent
	}

	return s
}

// partsToTree transforms a splitted JSON pointer to a tree
// The traversal is iterative to avoid unbounded recursion: depth would otherwise equal
// the number of pointer segments, which an attacker controls through the directive value
func partsToTree(t _type, parts []string, root *node, params *httpsfv.Params) {
	n := root
	for _, part := range parts {
		var child *node
		for _, c := range n.children {
			if c.path == part {
				child = c
				break
			}
		}

		if child == nil {
			child = &node{path: part, parent: n}
			n.children = append(n.children, child)
		}

		switch t {
		case preload:
			child.preload = true
			child.preloadParams = append(child.preloadParams, params)
		case fields:
			child.fields = true
			child.fieldsParams = append(child.fieldsParams, params)
		}

		n = child
	}
}

// hasChildren checks if the node has at least a child of the given type
func (n *node) hasChildren(t _type) bool {
	for _, c := range n.children {
		if t == preload && c.preload {
			return true
		}
		if t == fields && c.fields {
			return true
		}
	}

	return false
}

// httpList transforms the node in an HTTP Structured Field List
func (n *node) httpList(t _type, prefix string) httpsfv.List {
	if len(n.children) == 0 {
		if prefix == "" {
			return httpsfv.List{}
		}

		var list httpsfv.List
		switch t {
		case preload:
			for _, params := range n.preloadParams {
				list = append(list, httpsfv.Item{Value: prefix, Params: params})
			}
		case fields:
			for _, params := range n.fieldsParams {
				list = append(list, httpsfv.Item{Value: prefix, Params: params})
			}
		}

		return list
	}

	var list httpsfv.List
	for _, c := range n.children {
		if (t == preload && !c.preload) || (t == fields && !c.fields) {
			continue
		}

		list = append(list, c.httpList(t, prefix+"/"+c.path)...)
	}

	return list
}
