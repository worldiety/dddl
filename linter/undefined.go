package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type UndeclaredUsageInData struct {
	hint
	Parent            *parser.Data
	Name              *parser.Ident
	ExpectedQualifier parser.Qualifier
}

type UndeclaredUsageInWorkflow struct {
	hint
	Parent            *parser.Workflow
	Name              *parser.Ident
	ExpectedQualifier parser.Qualifier
}

// CheckUndefined searches for all Identifiers and checks
// if there are workflows or data types for them.
func CheckUndefined(root parser.Node) []Hint {
	ws := parser.WorkspaceOf(root)
	if ws == nil {
		return nil
	}

	dedupTableData := map[string]*UndeclaredUsageInData{}
	dedupTableWorkflow := map[string]*UndeclaredUsageInWorkflow{}

	err := parser.Walk(root, func(n parser.Node) error {
		if name, ok := n.(*parser.Ident); ok {
			if name.IsUniverse() {
				return nil
			}

			if parent := parser.DataOf(name); parent != nil {
				expected, declared := ws.Resolve(name)
				if !declared {
					hint := dedupTableData[expected.String()]
					if hint == nil {
						dedupTableData[expected.String()] = &UndeclaredUsageInData{
							Parent:            parent,
							Name:              name,
							ExpectedQualifier: expected,
						}
					}
				}
			}

			if parent := parser.WorkflowOf(name); parent != nil {
				expected, declared := ws.Resolve(name)
				if !declared {
					hint := dedupTableWorkflow[expected.String()]
					if hint == nil {
						dedupTableWorkflow[expected.String()] = &UndeclaredUsageInWorkflow{
							Parent:            parent,
							Name:              name,
							ExpectedQualifier: expected,
						}
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("cannot happen: %w", err))
	}

	var res []Hint
	keys := maps.Keys(dedupTableWorkflow)
	slices.Sort(keys)
	for _, key := range keys {
		res = append(res, dedupTableWorkflow[key])
	}

	keys = maps.Keys(dedupTableData)
	slices.Sort(keys)
	for _, key := range keys {
		res = append(res, dedupTableData[key])
	}

	return res
}
