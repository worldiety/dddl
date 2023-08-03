package linter

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type AmbiguousDeclaration struct {
	hint
	Declarations []*parser.TypeDefinition
}

// CheckAmbiguous validates that defined terms (e.g. Context, Workflow and Data)
// are all unique.
func CheckAmbiguous(r *resolver.Resolver) []Hint {

	allDefs := map[string][]*parser.TypeDefinition{}
	for _, context := range r.Contexts() {
		for _, fragment := range context.Fragments {
			for _, definition := range fragment.Definitions {
				qname := resolver.NewQualifiedNameFromNamedType(definition.Type)
				list := allDefs[qname.String()]
				list = append(list, definition)
				allDefs[qname.String()] = list
			}
		}
	}

	var res []Hint
	keys := maps.Keys(allDefs)
	slices.Sort(keys)
	for _, key := range keys {
		list := allDefs[key]
		if len(list) > 1 {
			res = append(res, &AmbiguousDeclaration{Declarations: list})
		}
	}

	return res
}
