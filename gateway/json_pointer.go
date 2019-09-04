package gateway

type jsonPointersTree struct {
	key    string // Root is zero valued
	values []*jsonPointersTree
}
