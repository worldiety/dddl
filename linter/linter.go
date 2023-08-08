package linter

import (
	"github.com/worldiety/dddl/resolver"
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
func Lint(r *resolver.Resolver) []Hint {
	var res []Hint
	res = append(res, CheckLiteralDefinitions(r)...)
	res = append(res, CheckUndefined(r)...)
	res = append(res, CheckAmbiguous(r)...)
	res = append(res, CheckAssignees(r)...)
	res = append(res, CheckNoContext(r)...)

	return res
}
