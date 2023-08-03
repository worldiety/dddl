package puml

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
)

func Type(r *resolver.Resolver, data *parser.Type, flags RFlags) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	iface := plantuml.NewInterface(data.Name.Value)
	diag.Add(iface)
	if flags.MainType == data {
		note := ""
		if data.Basetype != nil {
			note += "Neuer eigenständiger Datentyp\nmit dem Basistyp " + typeDeclToLinkStr(r, data.Basetype)
		} else {
			note = "Nicht näher bestimmter Datentyp."
		}

		iface.NoteRight(plantuml.NewNote(note))
	}

	if data.Basetype != nil {

		defs := r.ResolveLocalQualifier(data.Basetype.Name)
		for _, def := range defs {
			diag.Add(RenderNamedType(r, def.Type, flags).Renderables...)
			diag.Add(&plantuml.Association{
				Owner: iface.Name(),
				Child: def.Type.GetName().Value,
				Type:  plantuml.AssocExtension,
			})
		}

		addUniverse(diag, iface.Name(), parser.NewIdentWithParent(data, data.Basetype.Name.String()))

		insertTypeParams(r, iface.Name(), diag, data.Basetype, flags)
	}

	return diag
}

func addUniverse(diag *plantuml.Diagram, typeName string, ident *parser.Ident) {
	if !ident.IsUniverse() {
		return
	}

	diag.Add(&plantuml.Association{
		Child: typeName,
		Owner: ident.Value,
		Type:  plantuml.AssocExtension,
	})
}
