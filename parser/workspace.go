package parser

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Workspace is nothing to be parsed, but can be used to aggregate multiple documents into something bigger and
// apply the linters on that.
type Workspace struct {
	node
	Documents map[string]*Doc
	Error     error
}

// Docs returns a stable sorted list of Documents.
func (n *Workspace) Docs() []*Doc {
	var res []*Doc
	keys := maps.Keys(n.Documents)
	slices.Sort(keys)
	for _, key := range keys {
		res = append(res, n.Documents[key])
	}

	return res
}

// Children returns a stable sorted list of Documents.
func (n *Workspace) Children() []Node {
	var res []Node

	for _, doc := range n.Docs() {
		res = append(res, doc)
	}

	return res
}
