package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

var regexAssignee = regexp.MustCompile(`@\w+:?\s[^.?!\n]+`)
var regexName = regexp.MustCompile(`@\w+:?`)

type AssignedTask struct {
	Assignee string
	Task     string
}

func ParseAssignees(str string) []AssignedTask {
	var res []AssignedTask
	matches := regexAssignee.FindAllString(str, -1)
	for _, match := range matches {
		names := regexName.FindAllString(match, 1)
		if len(names) > 0 {
			name := names[0][1:] // cut of @
			if strings.HasSuffix(name, ":") {
				name = name[:len(name)-1]
			}
			res = append(res, AssignedTask{
				Assignee: name,
				Task:     strings.TrimSpace(match[len(names[0]):]),
			})
		}
	}

	return res
}

type AssignedDefinition struct {
	hint
	Task AssignedTask
	Def  *parser.TypeDefinition
}

type AssignedTasks struct {
	hint
	Assignee   string
	Categories map[string][]AssignedDefinition
}

// CheckAssignees inspects all definitions and Todos for assigned tasks.
func CheckAssignees(r *resolver.Resolver) []Hint {
	var res []Hint
	clustered := map[string]*AssignedTasks{}
	getter := func(task AssignedTask) *AssignedTasks {
		tmp := clustered[task.Assignee]
		if tmp == nil {
			tmp = &AssignedTasks{Assignee: task.Assignee, Categories: map[string][]AssignedDefinition{}}
			clustered[task.Assignee] = tmp
		}

		return tmp
	}

	for _, context := range r.Contexts() {

		for _, fragment := range context.Fragments {

			if ctxFragmentTypeDef, ok := fragment.Parent().(*parser.TypeDefinition); ok {
				if ctxFragmentTypeDef.Description != nil {

					for _, task := range ParseAssignees(ctxFragmentTypeDef.Description.Value) {
						tmp := getter(task)
						typename := fmt.Sprintf("%T", ctxFragmentTypeDef.Type)
						list := tmp.Categories[typename]
						list = append(list, AssignedDefinition{
							Task: task,
							Def:  ctxFragmentTypeDef,
						})
						tmp.Categories[typename] = list
					}
				}
			}

			for _, definition := range fragment.Definitions {
				if definition.Description != nil {
					for _, task := range ParseAssignees(definition.Description.Value) {
						tmp := getter(task)
						typename := fmt.Sprintf("%T", definition.Type)
						list := tmp.Categories[typename]
						list = append(list, AssignedDefinition{
							Task: task,
							Def:  definition,
						})
						tmp.Categories[typename] = list
					}

				}

			}
		}
	}

	keys := maps.Keys(clustered)
	slices.Sort(keys)
	for _, key := range keys {
		res = append(res, clustered[key])
	}

	return res
}
