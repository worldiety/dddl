package linter

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
	"strings"
)

type TypeDefinitionNotDescribed struct {
	hint
	Def          *parser.TypeDefinition
	HasOpenTasks bool
}

// CheckLiteralDefinitions inspects the "Definition" literals for types.
// Every parser.TypeDecl should have one, otherwise
// the ubiquitous language is incomplete.
func CheckLiteralDefinitions(r *resolver.Resolver) []Hint {
	var res []Hint

	for _, context := range r.Contexts() {
		for _, fragment := range context.Fragments {
			for _, definition := range fragment.Definitions {
				if definition.Description == nil {
					res = append(res, &TypeDefinitionNotDescribed{Def: definition})
					continue
				}

				text := strings.TrimSpace(strings.ToLower(definition.Description.Value))
				if text == "" {
					res = append(res, &TypeDefinitionNotDescribed{Def: definition})
					continue
				}

				if strings.Contains(text, "todo") {
					res = append(res, &TypeDefinitionNotDescribed{Def: definition, HasOpenTasks: true})
					continue
				}
			}
		}
	}

	return res
}
