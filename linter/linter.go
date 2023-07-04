package linter

import (
	"github.com/worldiety/dddl/parser"
)

type hint struct {
}

func (h *hint) Hint() bool {
	return true
}

type Hint interface {
	Hint() bool
}

// Lint applies all available linters.
func Lint(root parser.Node) []Hint {
	var res []Hint
	res = append(res, CheckToDos(root)...)
	res = append(res, CheckLiteralDefinitions(root)...)
	res = append(res, CheckUndefined(root)...)
	res = append(res, CheckAmbiguous(root)...)
	res = append(res, CheckAssignees(root)...)

	return res
}
