package vulcain

import (
	"strings"

	"github.com/dunglas/httpsfv"
)

type node struct {
	preload       bool
	preloadParams []*httpsfv.Params
	fields        bool
	fieldsParams  []*httpsfv.Params
	path          string
	parent        *node
	children      []*node
}

// Type is the type of operation to apply, can be Preload or Fields
type Type int

const (
	// Preload is a preloading action through query parameters or headers
	Preload Type = iota
	// Fields is a filtering action through query parameters or headers
	Fields
)

func (n *node) importPointers(t Type, pointers httpsfv.List) {
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

func partsToTree(t Type, parts []string, root *node, params *httpsfv.Params) {
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
	case Preload:
		child.preload = true
		child.preloadParams = append(child.preloadParams, params)
	case Fields:
		child.fields = true
		child.fieldsParams = append(child.fieldsParams, params)
	}

	partsToTree(t, parts[1:], child, params)
}

func (n *node) hasChildren(t Type) bool {
	for _, c := range n.children {
		if t == Preload && c.preload {
			return true
		}
		if t == Fields && c.fields {
			return true
		}
	}

	return false
}

func (n *node) httpList(t Type, prefix string) httpsfv.List {
	if len(n.children) == 0 {
		if prefix == "" {
			return httpsfv.List{}
		}

		var list httpsfv.List
		switch t {
		case Preload:
			for _, params := range n.preloadParams {
				list = append(list, httpsfv.Item{Value: prefix, Params: params})
			}
		case Fields:
			for _, params := range n.fieldsParams {
				list = append(list, httpsfv.Item{Value: prefix, Params: params})
			}
		}

		return list
	}

	var list httpsfv.List
	for _, c := range n.children {
		if (t == Preload && !c.preload) || (t == Fields && !c.fields) {
			continue
		}

		list = append(list, c.httpList(t, prefix+"/"+c.path)...)
	}

	return list
}
