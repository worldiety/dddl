package puml

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
)

func Alias(r *resolver.Resolver, data *parser.Alias, flags RFlags) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	iface := plantuml.NewInterface(data.Name.Value)
	diag.Add(iface)
	if flags.MainType == data {
		note := ""
		if data.BaseType != nil {
			note += "Alternativer Name für " + typeDeclToLinkStr(r, data.BaseType)
		} else {
			note = "Nicht näher bestimmter Datentyp."
		}

		iface.NoteRight(plantuml.NewNote(note))
	}

	if data.BaseType != nil {

		defs := r.ResolveLocalQualifier(data.BaseType.Name)
		for _, def := range defs {
			diag.Add(RenderNamedType(r, def.Type, flags).Renderables...)
			diag.Add(&plantuml.Association{
				Child: iface.Name(),
				Owner: def.Type.GetName().Value,
				Type:  plantuml.AssocExtension,
			})
		}

		addUniverse(diag, iface.Name(), parser.NewIdentWithParent(data, data.BaseType.Name.String()))

		insertTypeParams(r, iface.Name(), diag, data.BaseType, flags)
	}

	return diag
}
