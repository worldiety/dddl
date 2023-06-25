package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
)

// CheckLiteralDefinitions inspects the "Definition" literals for types.
// Every parser.TypeDefinition should have one, otherwise
// the ubiquitous language is incomplete.
func CheckLiteralDefinitions(root parser.Node) []Hint {
	var res []Hint
	err := parser.Walk(root, func(n parser.Node) error {
		switch n := n.(type) {
		case *parser.Context:
			if n.Definition.Empty() {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.Definition,
					Message:     "Die Beschreibung des Kontexts %s fehlt.",
				})
			}
		case *parser.Workflow:
			if n.Definition.Empty() {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.Definition,
					Message:     "Die Beschreibung des Arbeitsablaufs %s fehlt.",
				})
			}

			if n.Definition.NeedsRevise() {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.Definition,
					Message:     "Die Beschreibung des Arbeitsablaufs %s enthält noch ungeklärte Fragen.",
				})
			}

		case *parser.Data:
			if n.Definition.Empty() {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.Definition,
					Message:     "Die Beschreibung zu den Daten %s fehlt.",
				})
			}

			if n.Definition.NeedsRevise() {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.Definition,
					Message:     "Die Beschreibung zu den Daten %s enthält noch ungeklärte Fragen.",
				})
			}
		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("cannot happen: %w", err))
	}

	return res
}
