package gateway

import (
	"strings"
)

type node struct {
	preload  bool
	fields   bool
	path     string
	children []*node
}

type Type int

const (
	Preload Type = iota
	// False is a json false boolean
	Fields
)

func newPointersTree(p bool, f bool) *node {
	return &node{preload: p, fields: f}
}

func (n *node) importPointers(t Type, pointers []string) {
	for _, pointer := range pointers {
		pointer = strings.Trim(pointer, "/")
		if pointer == "" {
			continue
		}

		parts := strings.Split(pointer, "/")
		partsToTree(t, parts, n)
	}
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
		root.children = append(root.children, child)
		// TODO: sort by path
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
