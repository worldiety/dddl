package linter

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

type InvalidAnnotation struct {
	hint
	TypeDef *parser.TypeDefinition
	Error   error
}

// CheckAnnotations validates that only allowed annotations are used.
func CheckAnnotations(r *resolver.Resolver) []Hint {
	var res []Hint

	for _, context := range r.Contexts() {
		for _, fragment := range context.Fragments {
			for _, definition := range fragment.Definitions {

				if err := definition.ExpectOnlyOf("external", "Fremdsystem", "error", "Fehler", "Ereignis", "event"); err != nil {
					res = append(res, &InvalidAnnotation{
						TypeDef: definition,
						Error:   err,
					})
				}

				if _, err := parser.ParseEventAnnotation(definition); err != nil {
					res = append(res, &InvalidAnnotation{
						TypeDef: definition,
						Error:   err,
					})
				}

				if _, err := parser.ParseErrorAnnotation(definition); err != nil {
					res = append(res, &InvalidAnnotation{
						TypeDef: definition,
						Error:   err,
					})
				}

				if _, err := parser.ParseExternalSystemAnnotation(definition); err != nil {
					res = append(res, &InvalidAnnotation{
						TypeDef: definition,
						Error:   err,
					})
				}
			}
		}
	}

	return res
}
