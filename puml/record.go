package puml

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
)

func Record(r *resolver.Resolver, data *parser.Struct, flags RFlags) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	class := ClassFromRecord(r, data)
	if data == flags.MainType {
		class.NoteRight(plantuml.NewNote(record2Str(data)))
	}
	diag.Add(class)

	if flags.Depth <= 0 {
		return diag
	}

	for _, field := range data.Fields {
		defs := r.ResolveLocalQualifier(field.TypeDecl.Name)
		for _, def := range defs {
			fields := RenderNamedType(r, def.Type, flags).Renderables
			diag.Add(fields...)
			diag.Add(&plantuml.Association{
				Owner:            class.Name(),
				OwnerCardinality: "1",
				Child:            def.Type.GetName().Value,
				Type:             plantuml.AssocAggregation,
			})
		}

		insertTypeParams(r, class.Name(), diag, field.TypeDecl, flags)
	}

	return diag
}

func insertTypeParams(r *resolver.Resolver, ownername string, diag *plantuml.Diagram, field *parser.TypeDeclaration, flags RFlags) {
	if len(field.Params) == 0 {
		return
	}

	ident := parser.UniverseName(field.Name.String())
	mult := ident.IsList() || ident.IsMap()
	for _, param := range field.Params {
		defs := r.ResolveLocalQualifier(param.Name)
		for _, def := range defs {
			diag.Add(RenderNamedType(r, def.Type, flags).Renderables...)
			if mult {
				diag.Add(&plantuml.Association{
					Owner:            ownername,
					OwnerCardinality: "*",
					Child:            def.Type.GetName().Value,
					Type:             plantuml.AssocAggregation,
				})
			} else {
				diag.Add(&plantuml.Association{
					Owner:            ownername,
					OwnerCardinality: "1",
					Child:            def.Type.GetName().Value,
					Type:             plantuml.AssocAggregation,
				})
			}
		}
	}
}

func ClassFromRecord(r *resolver.Resolver, data *parser.Struct) *plantuml.Class {
	c := plantuml.NewClass(data.Name.Value)

	evtA := parser.FindAnnotation[*parser.EventAnnotation](data.Parent().(*parser.TypeDefinition))
	if evtA != nil {
		c.Stereotypes = append(c.Stereotypes, "Domänenereignis")
		if evtA.In {
			c.Stereotypes = append(c.Stereotypes, "eingehendes Ereignis")
		}

		if evtA.Out {
			c.Stereotypes = append(c.Stereotypes, "ausgehendes Ereignis")
		}
	}

	if a := parser.FindAnnotation[*parser.RoleAnnotation](data); a != nil {
		c.Stereotypes = append(c.Stereotypes, "Nutzer-Rolle")
	}

	evtErr := parser.FindAnnotation[*parser.ErrorAnnotation](data.Parent().(*parser.TypeDefinition))
	if evtErr != nil {
		c.Stereotypes = append(c.Stereotypes, "Fehler")
	}

	for _, f := range data.Fields {
		alias := ""
		if f.Alias != nil {
			alias = f.Alias.Value + ": "
		}
		c.AddAttrs(plantuml.Attr{
			Visibility: plantuml.Public,
			Name:       alias + typeDeclToLinkStr(r, f.TypeDecl),
		})

	}

	return c
}

func ClassFromChoice(r *resolver.Resolver, data *parser.Choice) *plantuml.Class {
	c := plantuml.NewClass(data.Name.Value)
	for _, declaration := range data.Choices {
		c.AddAttrs(plantuml.Attr{
			Visibility: plantuml.Public,
			Name:       typeDeclToLinkStr(r, declaration),
		})
	}

	return c
}
