package linter

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

type DeclaredWithoutContext struct {
	hint
	Parent  *parser.Doc
	TypeDef *parser.TypeDefinition
}

// CheckNoContext yells about all declarations within the anonymous space.
func CheckNoContext(r *resolver.Resolver) []Hint {
	var res []Hint
	for _, doc := range r.Workspace().Documents {
		for _, definition := range doc.Types {
			if _, isCtx := definition.Type.(*parser.Context); isCtx {
				continue
			}

			res = append(res, &DeclaredWithoutContext{
				Parent:  doc,
				TypeDef: definition,
			})

		}
	}

	return res
}
