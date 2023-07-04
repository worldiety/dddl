package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
)

type ToDoData struct {
	hint
	Parent *parser.Data
}

type ToDoWorkflow struct {
	hint
	Parent *parser.Workflow
}

type ToDoContext struct {
	hint
	Parent *parser.Context
}

// CheckToDos collects all relevant [parser.ToDo] entries.
func CheckToDos(root parser.Node) []Hint {
	var res []Hint
	err := parser.Walk(root, func(n parser.Node) error {
		if todo, ok := n.(*parser.ToDo); ok {
			if ctx, ok := todo.Parent().(*parser.Context); ok {
				res = append(res, &ToDoContext{Parent: ctx})
			}

			if data := parser.DataOf(n); data != nil {
				res = append(res, &ToDoData{Parent: data})
			}

			if wf := parser.WorkflowOf(n); wf != nil {
				res = append(res, &ToDoWorkflow{Parent: wf})
			}
		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("unreachable: %w", err))
	}

	return res
}
