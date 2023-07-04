package lsp

import (
	"fmt"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/lsp/protocol"
	"github.com/worldiety/dddl/parser"
	"path/filepath"
)

func newDiag(n parser.Node, msg string) protocol.Diagnostic {
	pos := n.Position()
	end := n.EndPosition()

	return protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{
				// Subtract 1 since dyml has 1 based lines and columns, but LSP wants 0 based
				Line:      uint32(pos.Line) - 1,
				Character: uint32(pos.Column) - 1,
			},
			// we don't know the length, so just always pick the next 3 chars
			End: protocol.Position{
				Line:      uint32(end.Line) - 1,
				Character: uint32(end.Column) - 1,
			},
		},
		Severity: protocol.SeverityWarning,
		Message:  msg,
	}
}

func renderLintTexts(matchFile protocol.DocumentURI, hints []linter.Hint) []protocol.Diagnostic {
	var res []protocol.Diagnostic
	matches := func(n parser.Node) bool {
		hintFname := filepath.Base(n.Position().Filename)
		baseFname := filepath.Base(string(matchFile)) // TODO assumption not correct for files in distinct folders
		return baseFname == hintFname
	}

	for _, hint := range hints {
		switch h := hint.(type) {
		case *linter.AmbiguousDeclaration:
			for _, declaration := range h.Declarations {
				if matches(declaration.DeclaredName()) {
					res = append(res, newDiag(h.Declarations[0].DeclaredName(), "Dieser Bezeichner wurde bereits woanders deklariert"))
				}
			}
		case *linter.AssignedContext:
			if matches(h.Context) {
				res = append(res, newDiag(h.Context.Name, h.Task.Assignee+" hat eine offene Aufgabe: "+h.Task.Task))
			}
		case *linter.AssignedData:
			if matches(h.Data) {
				res = append(res, newDiag(h.Data.Name, h.Task.Assignee+" hat eine offene Aufgabe: "+h.Task.Task))
			}

		case *linter.AssignedTasks:
			// ignore, for linting, we use the inflated ones
		case *linter.AssignedWorkflow:
			if matches(h.Workflow) {
				res = append(res, newDiag(h.Workflow.Name, h.Task.Assignee+" hat eine offene Aufgabe: "+h.Task.Task))
			}
		case *linter.ContextHasQuestions:
			if matches(h.Context) {
				res = append(res, newDiag(h.Context.Name, "Der Kontext enthält noch offene Fragen."))
			}
		case *linter.ContextNotDescribed:
			if matches(h.Context) {
				res = append(res, newDiag(h.Context.Name, "Der Kontext wurde noch nicht beschrieben."))
			}
		case *linter.DataHasQuestions:
			if matches(h.Data) {
				res = append(res, newDiag(h.Data.Name, "Die Daten-Deklaration enthält noch offene Fragen."))
			}
		case *linter.DataNotDescribed:
			if matches(h.Data) {
				res = append(res, newDiag(h.Data.Name, "Die Daten-Deklaration hat keine Beschreibung."))
			}
		case *linter.ToDoContext:
			if matches(h.Parent) {
				res = append(res, newDiag(h.Parent.Name, "Die Kontext-Deklaration hat ein offenes TODO"))
			}
		case *linter.ToDoData:
			if matches(h.Parent) {
				res = append(res, newDiag(h.Parent.Name, "Die Daten-Deklaration hat ein offenes TODO."))
			}
		case *linter.ToDoWorkflow:
			if matches(h.Parent) {
				res = append(res, newDiag(h.Parent.Name, "Die Arbeitsablauf-Deklaration hat ein offenes TODO"))
			}
		case *linter.UndeclaredUsageInData:
			if matches(h.Parent) {
				res = append(res, newDiag(h.Parent.Name, "Die Daten-Deklaration verwendet den unbekannten Bezeichner "+h.Name.Value))
			}
		case *linter.UndeclaredUsageInWorkflow:
			if matches(h.Parent) {
				res = append(res, newDiag(h.Parent.Name, "Der Arbeitsablauf verwendet den unbekannten Bezeichner "+h.Name.Value))
			}
		case *linter.WorkflowHasQuestions:
			if matches(h.Workflow) {
				res = append(res, newDiag(h.Workflow.Name, "Der Arbeitsablauf enthält noch offene Fragen."))
			}

		case *linter.WorkflowNotDescribed:
			if matches(h.Workflow) {
				res = append(res, newDiag(h.Workflow.Name, "Der Arbeitsablauf ist noch nicht beschrieben."))
			}

		default:
			panic(fmt.Sprintf("implement lint support: %T", h))
		}
	}

	return res
}
