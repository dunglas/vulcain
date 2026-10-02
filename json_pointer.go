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
	sources       [2]*node
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
	if n.sources[preload] != nil {
		return n.sources[preload].String()
	}
	if n.sources[fields] != nil {
		return n.sources[fields].String()
	}

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
	source := n.source(t)
	if source == nil {
		return false
	}

	for _, c := range source.children {
		if t == preload && c.preload {
			return true
		}
		if t == fields && c.fields {
			return true
		}
	}

	return false
}

func (n *node) source(t _type) *node {
	if n.sources[preload] != nil || n.sources[fields] != nil {
		return n.sources[t]
	}

	return n
}

// child applies exact-over-wildcard precedence independently for each directive.
func (n *node) child(key []byte) *node {
	p := n.match(preload, key)
	f := n.match(fields, key)
	if p == f {
		return p
	}
	if p != nil && f == nil && !p.fields {
		return p
	}
	if f != nil && p == nil && !f.preload {
		return f
	}

	child := &node{sources: [2]*node{p, f}}
	if p != nil {
		child.preload = true
		child.preloadParams = p.preloadParams
	}
	if f != nil {
		child.fields = true
		child.fieldsParams = f.fieldsParams
	}

	return child
}

func (n *node) match(t _type, key []byte) *node {
	source := n.source(t)
	if source == nil {
		return nil
	}

	var wildcard *node
	for _, c := range source.children {
		if (t == preload && !c.preload) || (t == fields && !c.fields) {
			continue
		}

		if c.path == "*" {
			wildcard = c
			continue
		}

		if string(key) == unescape(c.path) {
			return c
		}
	}

	return wildcard
}

// httpList transforms the node in an HTTP Structured Field List
func (n *node) httpList(t _type, prefix string) httpsfv.List {
	source := n.source(t)
	if source == nil {
		return nil
	}
	if source != n {
		return source.httpList(t, prefix)
	}

	if !n.hasChildren(t) {
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
