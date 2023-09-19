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
				_ = definition
				//TODO
			}
		}
	}

	return res
}
