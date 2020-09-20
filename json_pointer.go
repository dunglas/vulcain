package vulcain

import (
	"strings"

	"github.com/dunglas/httpsfv"
)

// node represend the a node in a JSON document
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
	// preload is a preloading action through query parameters or headers
	preload _type = iota
	// fields is a filtering action through query parameters or headers
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

// String converts the tree as a JSON pointer
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
func partsToTree(t _type, parts []string, root *node, params *httpsfv.Params) {
	if len(parts) == 0 {
		return
	}

	var child *node
	for _, c := range root.children {
		if c.path == parts[0] {
			child = c
			break
		}
	}

	if child == nil {
		child = &node{}
		child.path = parts[0]
		child.parent = root
		root.children = append(root.children, child)
	}

	switch t {
	case preload:
		child.preload = true
		child.preloadParams = append(child.preloadParams, params)
	case fields:
		child.fields = true
		child.fieldsParams = append(child.fieldsParams, params)
	}

	partsToTree(t, parts[1:], child, params)
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

// httpList transforms the node an HTTP Structured Field List
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
