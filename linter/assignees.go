package linter

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
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

type AssignedContext struct {
	hint
	Task    AssignedTask
	Context *parser.Context // from to-do or definition
}

type AssignedData struct {
	hint
	Task AssignedTask
	Data *parser.Data // from to-do or definition
}

type AssignedWorkflow struct {
	hint
	Task     AssignedTask
	Workflow *parser.Workflow // from to-do or definition or nested to-do within the actual workflow
}

type AssignedTasks struct {
	hint
	Assignee  string
	Contexts  []AssignedContext
	Workflows []AssignedWorkflow
	Datas     []AssignedData
}

// CheckAssignees inspects all definitions and Todos for assigned tasks.
func CheckAssignees(root parser.Node) []Hint {
	var res []Hint
	clustered := map[string]*AssignedTasks{}
	getter := func(task AssignedTask) *AssignedTasks {
		tmp := clustered[task.Assignee]
		if tmp == nil {
			tmp = &AssignedTasks{Assignee: task.Assignee}
			clustered[task.Assignee] = tmp
		}

		return tmp
	}

	err := parser.Walk(root, func(n parser.Node) error {
		switch n := n.(type) {
		case *parser.Context:
			for _, task := range ParseAssignees(n.Definition.Value()) {
				tmp := getter(task)
				tmp.Contexts = append(tmp.Contexts, AssignedContext{Task: task, Context: n})
				res = append(res, &AssignedContext{Task: task, Context: n})
			}

			for _, task := range ParseAssignees(n.ToDo.Value()) {
				tmp := getter(task)
				tmp.Contexts = append(tmp.Contexts, AssignedContext{Task: task, Context: n})
				res = append(res, &AssignedContext{Task: task, Context: n})
			}

		case *parser.Workflow:
			for _, task := range ParseAssignees(n.Definition.Value()) {
				tmp := getter(task)
				tmp.Workflows = append(tmp.Workflows, AssignedWorkflow{Task: task, Workflow: n})
				res = append(res, &AssignedWorkflow{Task: task, Workflow: n})
			}

			for _, task := range ParseAssignees(n.ToDo.Value()) {
				tmp := getter(task)
				tmp.Workflows = append(tmp.Workflows, AssignedWorkflow{Task: task, Workflow: n})
				res = append(res, &AssignedWorkflow{Task: task, Workflow: n})
			}

			err := parser.Walk(n.Block, func(node parser.Node) error {
				if todo, ok := node.(*parser.ToDo); ok {
					for _, task := range ParseAssignees(todo.Value()) {
						tmp := getter(task)
						tmp.Workflows = append(tmp.Workflows, AssignedWorkflow{Task: task, Workflow: n})
						res = append(res, &AssignedWorkflow{Task: task, Workflow: n})
					}
				}

				return nil
			})

			if err != nil {
				return err
			}

		case *parser.Data:
			for _, task := range ParseAssignees(n.Definition.Value()) {
				tmp := getter(task)
				tmp.Datas = append(tmp.Datas, AssignedData{Task: task, Data: n})
				res = append(res, &AssignedData{Task: task, Data: n})
			}

			for _, task := range ParseAssignees(n.ToDo.Value()) {
				tmp := getter(task)
				tmp.Datas = append(tmp.Datas, AssignedData{Task: task, Data: n})
				res = append(res, &AssignedData{Task: task, Data: n})
			}
		}

		return nil
	})

	if err != nil {
		panic(fmt.Errorf("cannot happen: %w", err))
	}

	keys := maps.Keys(clustered)
	slices.Sort(keys)
	for _, key := range keys {
		res = append(res, clustered[key])
	}
	
	return res
}
