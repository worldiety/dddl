package html

import (
	"fmt"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"html/template"
)

func transformLintHints(r *resolver.Resolver, hints []linter.Hint, model PreviewModel) PreviewModel {

	for _, hint := range hints {
		var entry string
		switch h := hint.(type) {
		case *linter.AmbiguousDeclaration:
			entry = fmt.Sprintf(`Die Definition %s ist mehrfach erfolgt, aber nur eine ist zulässig:<ul class="list-disc">`, h.Declarations[0].Type.GetName().Value)
			for _, d := range h.Declarations {
				entry += "<li>" + filePos(d.Type.GetName()) + "</li>"
			}
			entry += "</ul>"
		case *linter.AssignedDefinition:
			typeDefName := CategoryName(h.Def.Type)
			entry = fmt.Sprintf(`%s %s enthält eine offene Aufgabe für %s: %s`, typeDefName, href(r, h.Def.Type), h.Task.Assignee, h.Task.Task)
		case *linter.AssignedTasks:
			nt := NamedTasks{Name: h.Assignee}
			keys := maps.Keys(h.Categories)
			slices.Sort(keys)
			for _, key := range keys {
				cat := h.Categories[key]
				for _, definition := range cat {
					nt.Tasks = append(nt.Tasks, template.HTML(fmt.Sprintf("in %s %s: %s", CategoryNameStr(key), href(r, definition.Def.Type), definition.Task.Task)))
				}
			}

			model.NamedTasks = append(model.NamedTasks, nt)
			// jump over, this will be shown in its own section template (NamedTasks)
			continue
		case *linter.FirstUndeclaredTypeDeclInNamedType:
			entry = fmt.Sprintf(`Die Definition %s verwendet den undefinierten Bezeichner %s`, href(r, h.Parent.Type), puml.TypeDeclToStr(h.TypeDecl))
		case *linter.UndeclaredTypeDeclInNamedType:
			// we only show the first occurence in our report
			continue
		case *linter.TypeDefinitionNotDescribed:
			typeDefName := CategoryName(h.Def.Type)
			entry = fmt.Sprintf(`%s %s hat noch keine vollständige Beschreibung.`, typeDefName, href(r, h.Def.Type))
		case *linter.DeclaredWithoutContext:
			typeDefName := CategoryName(h.TypeDef.Type)
			entry = fmt.Sprintf(`%s %s ist fachlich nicht zugeordnet und muss in einem Bounded Context definiert werden.`, typeDefName, href(r, h.TypeDef.Type))
		default:
			entry = fmt.Sprintf("implement: %T", hint)
		}

		model.Hints = append(model.Hints, template.HTML(entry))

	}

	return model
}

func CategoryName(namedType parser.NamedType) string {
	return CategoryNameStr(fmt.Sprintf("%T", namedType))
}

func CategoryNameStr(namedType string) string {
	typeDefName := namedType
	switch namedType {
	case "*parser.Context":
		typeDefName = "Kontext"
	case "*parser.Struct":
		typeDefName = "Datentyp"
	case "*parser.Type":
		typeDefName = "Basistyp"
	case "*parser.Choice":
		typeDefName = "Auswahltyp"
	case "*parser.Alias":
		typeDefName = "Alias"
	case "*parser.Function":
		typeDefName = "Aufgabe"
	case "*parser.Aggregate":
		typeDefName = "Aggregat"
	}

	return typeDefName
}

func href(r *resolver.Resolver, n parser.NamedType) string {
	qname := resolver.NewQualifiedNameFromNamedType(n)

	return fmt.Sprintf(`<a class="text-green-600" href="#%s">%s</a>`, qname.String(), qname.Name())
}

func filePos(n parser.Node) string {
	return fmt.Sprintf("%s:%d", n.Position().Filename, n.Position().Line)
}
