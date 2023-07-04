package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type AmbiguousDeclaration struct {
	hint
	Declarations []parser.Declaration
}

// CheckAmbiguous validates that defined terms (e.g. Context, Workflow and Data)
// are all unique.
func CheckAmbiguous(root parser.Node) []Hint {
	ws := parser.WorkspaceOf(root)
	if ws == nil {
		return nil
	}

	allDefs := map[string][]parser.Declaration{}
	err := parser.Walk(root, func(n parser.Node) error {
		if decl, ok := n.(parser.Declaration); ok {
			if _, isCtx := decl.(*parser.Context); isCtx {
				return nil
			}

			q, ok := ws.Resolve(decl.DeclaredName())
			if ok {
				list := allDefs[q.String()]
				list = append(list, decl)
				allDefs[q.String()] = list
			}
		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("cannot happen: %w", err))
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
