package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
)

type ContextNotDescribed struct {
	hint
	Context *parser.Context
}

type ContextHasQuestions struct {
	hint
	Context *parser.Context
}

type WorkflowHasQuestions struct {
	hint
	Workflow *parser.Workflow
}

type WorkflowNotDescribed struct {
	hint
	Workflow *parser.Workflow
}

type DataNotDescribed struct {
	hint
	Data *parser.Data
}

type DataHasQuestions struct {
	hint
	Data *parser.Data
}

// CheckLiteralDefinitions inspects the "Definition" literals for types.
// Every parser.TypeDecl should have one, otherwise
// the ubiquitous language is incomplete.
func CheckLiteralDefinitions(root parser.Node) []Hint {
	var res []Hint
	err := parser.Walk(root, func(n parser.Node) error {
		switch n := n.(type) {
		case *parser.Context:
			if n.Definition.Empty() {
				res = append(res, &ContextNotDescribed{Context: n})
			}

			if n.Definition.NeedsRevise() {
				res = append(res, &ContextHasQuestions{Context: n})
			}
		case *parser.Workflow:
			if n.Definition.Empty() {
				res = append(res, &WorkflowNotDescribed{Workflow: n})
			}

			if n.Definition.NeedsRevise() {
				res = append(res, &WorkflowHasQuestions{Workflow: n})
			}

		case *parser.Data:
			if n.Definition.Empty() {
				res = append(res, &DataNotDescribed{Data: n})
			}

			if n.Definition.NeedsRevise() {
				res = append(res, &DataHasQuestions{Data: n})
			}
		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("cannot happen: %w", err))
	}

	return res
}
