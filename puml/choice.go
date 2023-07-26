package puml

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
)

func Choice(r *resolver.Resolver, data *parser.Choice, flags RFlags) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	iface := plantuml.NewInterface(data.Name.Value)
	diag.Add(iface)
	if flags.MainType == data {
		iface.NoteRight(plantuml.NewNote(choice2Str(data)))
	}
	for _, choice := range data.Choices {
		choiceName := TypeDeclToStr(choice)

		defs := r.ResolveLocalQualifier(choice.Name)
		if len(defs) == 0 {
			diag.Add(plantuml.NewClass(choiceName).Extends(data.Name.Value))
		} else {
			for _, def := range defs {
				diag.Add(RenderNamedType(r, def.Type, flags).Renderables...)
				diag.Add(&plantuml.Association{
					Owner: iface.Name(),
					Child: def.Type.GetName().Value,
					Type:  plantuml.AssocExtension,
				})
			}

			insertTypeParams(r, iface.Name(), diag, choice, flags)
		}

	}

	return diag
}
