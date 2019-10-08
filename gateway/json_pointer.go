package gateway

import (
	"strings"
)

type node struct {
	preload  bool
	fields   bool
	path     string
	parent   *node
	children []*node
}

// Type is the type of operation to apply, can be Preload or Fields
type Type int

const (
	// Preload is a preloading action through query parameters or headers
	Preload Type = iota
	// Fields is a filtering action through query parameters or headers
	Fields
)

func (n *node) importPointers(t Type, pointers []string) {
	for _, pointer := range pointers {
		pointer = strings.Trim(pointer, "/")
		if pointer != "" {
			partsToTree(t, strings.Split(pointer, "/"), n)
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

func partsToTree(t Type, parts []string, root *node) {
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
		break

	case Fields:
		child.fields = true
		break
	}

	partsToTree(t, parts[1:], child)
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

func (n *node) strings(t Type, prefix string) []string {
	if len(n.children) == 0 {
		if prefix == "" {
			return []string{"/"}
		}

		return []string{prefix}
	}

	var pointers []string
	for _, c := range n.children {
		if t == Preload && !c.preload {
			continue
		}
		if t == Fields && !c.fields {
			continue
		}

		pointers = append(pointers, c.strings(t, prefix+"/"+c.path)...)
	}

	return pointers
}
