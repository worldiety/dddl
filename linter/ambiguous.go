package linter

import "github.com/worldiety/dddl/parser"

// CheckAmbiguous validates that defined terms (e.g. Context, Workflow and Data)
// are all unique.
func CheckAmbiguous(root parser.Node) []Hint {
	allDefs := map[string]parser.Node{}
	var res []Hint
	_ = parser.Walk(root, func(n parser.Node) error {
		var parent *parser.Ident
		var node parser.Node

		switch n := n.(type) {
		case *parser.Data:
			parent = n.Name
			node = n
		case *parser.Workflow:
			parent = n.Name
			node = n
		case *parser.Context:
			parent = n.Name
			node = n
		}

		if parent != nil {
			if other, ok := allDefs[parent.Value]; ok {
				res = append(res, Hint{
					ParentIdent: parent,
					Node:        other,
					Message:     "Der Begriff %s wurde mehrfach definiert.",
				})
			} else {
				allDefs[parent.Value] = node
			}
		}

		return nil
	})

	return res
}
