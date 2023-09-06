package lsp

import (
	_ "embed"
	"fmt"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/lsp/protocol"
	"github.com/worldiety/dddl/parser"
	"path/filepath"
)

//go:embed tips/help_missingdescription.md
var help_missingdescription string

//go:embed tips/help_tododescription.md
var help_tododescription string

//go:embed tips/help_nocontext.md
var help_nocontext string

//go:embed tips/help_undefined.md
var help_undefined string

//go:embed tips/help_invalidannotation.md
var help_invalidAnnotation string

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
				if matches(declaration.Type.GetName()) {
					res = append(res, newDiag(h.Declarations[0].Type.GetName(), "Dieser Bezeichner wurde bereits woanders deklariert: "+declaration.Pos.String()))
				}
			}
		case *linter.AssignedDefinition:
			if matches(h.Def) {
				res = append(res, newDiag(h.Def.Type.GetName(), h.Task.Assignee+" hat eine offene Aufgabe: "+h.Task.Task))
			}

		case *linter.AssignedTasks:
			// ignore, for linting, we use the inflated ones

		case *linter.TypeDefinitionNotDescribed:
			if matches(h.Def) {
				keyword := h.Def.Type.GetKeyword()
				typeName := h.Def.Type.GetName().Value
				if h.HasOpenTasks {
					res = append(res, newDiag(h.Def.Type.GetName(), fmt.Sprintf(help_tododescription, keyword, typeName)))
				} else {
					res = append(res, newDiag(h.Def.Type.GetName(), fmt.Sprintf(help_missingdescription, keyword, typeName)))
				}
			}

		case *linter.UndeclaredTypeDeclInNamedType:
			typeName := h.TypeDecl.Name.String()
			if matches(h.TypeDecl) {
				res = append(res, newDiag(h.TypeDecl.Name, fmt.Sprintf(help_undefined, typeName)))
			}

		case *linter.FirstUndeclaredTypeDeclInNamedType:
			continue
		case *linter.DeclaredWithoutContext:
			keyword := h.TypeDef.Type.GetKeyword()
			typeName := h.TypeDef.Type.GetName().Value
			if matches(h.TypeDef) {
				res = append(res, newDiag(h.TypeDef.Type.GetName(), fmt.Sprintf(help_nocontext, keyword, typeName)))
			}
		case *linter.InvalidAnnotation:
			typeName := h.TypeDef.Type.GetName().Value
			if matches(h.TypeDef) {
				res = append(res, newDiag(h.TypeDef.Type.GetName(), fmt.Sprintf(help_invalidAnnotation, typeName, h.Error.Error())))
			}
		default:
			panic(fmt.Sprintf("implement lint support: %T", h))
		}
	}

	return res
}
