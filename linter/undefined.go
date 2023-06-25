package linter

import (
	"github.com/worldiety/dddl/parser"
)

// CheckUndefined searches for all Identifiers and checks
// if there are workflows or data types for them.
func CheckUndefined(root parser.Node) []Hint {

	type Holder struct {
		Parent *parser.Ident
		Ident  *parser.Ident
	}

	var allIdents []Holder
	definedIdents := map[string]parser.Node{}

	_ = parser.Walk(root, func(n parser.Node) error {
		switch n := n.(type) {
		case *parser.Data, *parser.Workflow:
			var parentId *parser.Ident
			if n, ok := n.(*parser.Data); ok {
				parentId = n.Name
			}

			if n, ok := n.(*parser.Workflow); ok {
				parentId = n.Name
			}

			definedIdents[parentId.Name] = n

			_ = parser.Walk(n, func(p parser.Node) error {
				if usedId, ok := p.(*parser.Ident); ok {
					allIdents = append(allIdents, Holder{
						Parent: parentId,
						Ident:  usedId,
					})
				}

				return nil
			})
		}

		return nil
	})

	var res []Hint
	for _, holder := range allIdents {
		if _, ok := definedIdents[holder.Ident.Name]; !ok && !holder.Ident.IsUniverse() {
			res = append(res, Hint{
				ParentIdent: holder.Parent,
				Node:        holder.Ident,
				Message:     "%s verwendet den undefinierten Begriff " + holder.Ident.Name,
			})
		}
	}

	return res
}
