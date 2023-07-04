package editor

import (
	"fmt"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/parser"
	"html/template"
)

func lintOnly(ws *parser.Workspace, lint Linter, model EditorPreview) EditorPreview {

	//fmt.Printf("%#v", model)

	for _, hint := range lint(ws) {
		var entry string
		switch h := hint.(type) {
		case *linter.AmbiguousDeclaration:
			entry = fmt.Sprintf(`Die Deklaration %s ist mehrfach erfolgt, aber nur eine ist zulässig:<ul class="list-disc">`, h.Declarations[0].DeclaredName().Value)
			for _, d := range h.Declarations {
				entry += "<li>" + filePos(d.DeclaredName()) + "</li>"
			}
			entry += "</ul>"
		case *linter.AssignedContext:
			entry = fmt.Sprintf(`Der Kontext %s enthält eine offene Aufgabe für %s: %s`, href(ws, h.Context.Name), h.Task.Assignee, h.Task.Task)
		case *linter.AssignedData:
			entry = fmt.Sprintf(`Der Datentyp %s enthält eine offene Aufgabe für %s: %s`, href(ws, h.Data.Name), h.Task.Assignee, h.Task.Task)
		case *linter.AssignedTasks:
			nt := NamedTasks{Name: h.Assignee}
			for _, n := range h.Contexts {
				nt.Tasks = append(nt.Tasks, template.HTML(fmt.Sprintf("Kontext %s: %s", href(ws, n.Context.Name), n.Task.Task)))
			}

			for _, workflow := range h.Workflows {
				nt.Tasks = append(nt.Tasks, template.HTML(fmt.Sprintf("Arbeitsablauf %s: %s", href(ws, workflow.Workflow.Name), workflow.Task.Task)))
			}

			for _, n := range h.Datas {
				nt.Tasks = append(nt.Tasks, template.HTML(fmt.Sprintf("Daten %s: %s", href(ws, n.Data.Name), n.Task.Task)))
			}

			model.NamedTasks = append(model.NamedTasks, nt)

			// jump over, this will be shown in its own section
			continue
		case *linter.AssignedWorkflow:
			entry = fmt.Sprintf(`Der Arbeitsablauf %s enthält eine offene Aufgabe für %s: %s`, href(ws, h.Workflow.Name), h.Task.Assignee, h.Task.Task)
		case *linter.ContextHasQuestions:
			entry = fmt.Sprintf(`Der Kontext %s enthält ungeklärte Fragen.`, href(ws, h.Context.Name))
		case *linter.ContextNotDescribed:
			entry = fmt.Sprintf(`Der Kontext %s hat noch keine Beschreibung.`, href(ws, h.Context.Name))
		case *linter.DataHasQuestions:
			entry = fmt.Sprintf(`Der Datentyp %s enthält ungeklärte Fragen.`, href(ws, h.Data.Name))
		case *linter.DataNotDescribed:
			entry = fmt.Sprintf(`Der Datentyp %s hat noch keine Beschreibung.`, href(ws, h.Data.Name))
		case *linter.ToDoContext:
			entry = fmt.Sprintf(`Der Kontext %s enthält noch Aufgaben:<div>%s</div>`, href(ws, h.Parent.Name), markdown(h.Parent.ToDo.Value()))
		case *linter.ToDoData:
			entry = fmt.Sprintf(`Der Datentype %s enthält noch Aufgaben:<div>%s</div>`, href(ws, h.Parent.Name), markdown(h.Parent.ToDo.Value()))
		case *linter.ToDoWorkflow:
			entry = fmt.Sprintf(`Der Arbeitsablauf %s enthält noch Aufgaben:<div>%s</div>`, href(ws, h.Parent.Name), markdown(h.Parent.ToDo.Value()))
		case *linter.UndeclaredUsageInData:
			entry = fmt.Sprintf(`Der Datentyp %s verwendet den undeklarierten Bezeichner %s`, href(ws, h.Parent.Name), h.Name.Value)
		case *linter.UndeclaredUsageInWorkflow:
			entry = fmt.Sprintf(`Der Arbeitsablauf %s verwendet den undeklarierten Bezeichner %s`, href(ws, h.Parent.Name), h.Name.Value)
		case *linter.WorkflowHasQuestions:
			entry = fmt.Sprintf(`Der Arbeitsablauf %s enthält ungeklärte Fragen.`, href(ws, h.Workflow.Name))
		case *linter.WorkflowNotDescribed:
			entry = fmt.Sprintf(`Der Arbeitsablauf %s hat noch keine Beschreibung.`, href(ws, h.Workflow.Name))
		default:
			entry = fmt.Sprintf("implement: %v", hint)
		}

		model.Hints = append(model.Hints, template.HTML(entry))

	}

	return model
}

func href(ws *parser.Workspace, n *parser.Ident) string {
	q, _ := ws.Resolve(n)
	return fmt.Sprintf(`<a class="text-green-600" href="#%s">%s</a>`, q.String(), n.Value)
}

func filePos(n parser.Node) string {
	return fmt.Sprintf("%s:%d", n.Position().Filename, n.Position().Line)
}
