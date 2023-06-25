package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
)

// CheckToDos collects all relevant [parser.ToDo] entries.
func CheckToDos(root parser.Node) []Hint {
	var res []Hint
	err := parser.Walk(root, func(n parser.Node) error {
		switch n := n.(type) {
		case *parser.Context:
			if n.ToDo != nil {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.ToDo,
					Message:     n.ToDo.Text.Text,
				})
			}
		case *parser.Workflow:
			if n.ToDo != nil {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.ToDo,
					Message:     n.ToDo.Text.Text,
				})
			}

			// special check for stmts
			err := parser.Walk(n.Block, func(pn parser.Node) error {
				if todo, ok := pn.(*parser.ToDo); ok {
					res = append(res, Hint{
						ParentIdent: n.Name,
						Node:        todo,
						Message:     todo.Text.Text,
					})
				}

				return nil
			})

			if err != nil {
				panic(fmt.Errorf("unreachable: %w", err))
			}

		case *parser.Data:
			if n.ToDo != nil {
				res = append(res, Hint{
					ParentIdent: n.Name,
					Node:        n.ToDo,
					Message:     n.ToDo.Text.Text,
				})
			}

		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("unreachable: %w", err))
	}

	return res
}
