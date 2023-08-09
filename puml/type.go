package puml

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
)

func Type(r *resolver.Resolver, data *parser.Type, flags RFlags) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	iface := plantuml.NewClass(data.Name.Value)
	diag.Add(iface)
	//if flags.MainType == data {
	note := ""
	if data.Basetype != nil {
		note += "Neuer eigenständiger Datentyp\nmit dem Basistyp " + typeDeclToLinkStr(r, data.Basetype)
	} else {
		note = "Nicht näher bestimmter Datentyp."
	}

	iface.NoteBottom(plantuml.NewNote(note))
	//}

	if data.Basetype != nil {

		defs := r.ResolveLocalQualifier(data.Basetype.Name)
		for _, def := range defs {
			diag.Add(RenderNamedType(r, def.Type, flags).Renderables...)
			diag.Add(&plantuml.Association{
				Child: iface.Name(),
				Owner: def.Type.GetName().Value,
				Type:  plantuml.AssocExtension,
			})
		}

		if flags.Depth <= 0 {
			return diag
		}

		addUniverse(diag, iface.Name(), parser.UniverseName(data.Basetype.Name.String()))

		insertTypeParams(r, iface.Name(), diag, data.Basetype, flags)
	}

	return diag
}

func addUniverse(diag *plantuml.Diagram, typeName string, name parser.UniverseName) {
	if !name.IsUniverse() {
		return
	}

	diag.Add(&plantuml.Association{
		Child: typeName,
		Owner: string(name),
		Type:  plantuml.AssocExtension,
	})
	/*
		note := plantuml.NewNote(string(name))
		note.Dir = "right"
		note.Node = typeName
		diag.Add(note)*/
}
