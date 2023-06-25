package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"strings"
)

type Hint struct {
	ParentIdent *parser.Ident // the identifier of the parent, e.g. for linking
	Node        parser.Node   // the affected nearest node, e.g. for marking a line
	Message     string        // the message to display
}

func (h Hint) String(render func(ident *parser.Ident) string) string {
	if strings.Contains(h.Message, "%s") {
		return fmt.Sprintf(h.Message, render(h.ParentIdent))
	}

	return h.Message
}

// Lint applies all available linters.
func Lint(root parser.Node) []Hint {
	var res []Hint
	res = append(res, CheckToDos(root)...)
	res = append(res, CheckLiteralDefinitions(root)...)
	res = append(res, CheckUndefined(root)...)
	res = append(res, CheckAmbiguous(root)...)

	return Unique(res)
}

func Unique(hints []Hint) []Hint {
	var res []Hint
	tmp := map[string]struct{}{}
	for _, hint := range hints {
		key := hint.Message + hint.ParentIdent.Name
		if _, ok := tmp[key]; ok {
			continue
		}

		tmp[key] = struct{}{}
		res = append(res, hint)
	}

	return res
}
