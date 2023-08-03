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

	for _, field := range data.Fields {
		defs := r.ResolveLocalQualifier(field.Name)
		for _, def := range defs {
			diag.Add(RenderNamedType(r, def.Type, flags).Renderables...)
			diag.Add(&plantuml.Association{
				Owner:            class.Name(),
				OwnerCardinality: "1",
				Child:            def.Type.GetName().Value,
				Type:             plantuml.AssocAggregation,
			})
		}

		insertTypeParams(r, class.Name(), diag, field, flags)
	}
	return diag
}

func insertTypeParams(r *resolver.Resolver, ownername string, diag *plantuml.Diagram, field *parser.TypeDeclaration, flags RFlags) {
	if len(field.Params) == 0 {
		return
	}

	ident := parser.NewIdentWithParent(field, field.Name.String())
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
	for _, declaration := range data.Fields {
		c.AddAttrs(plantuml.Attr{
			Visibility: plantuml.Public,
			Name:       typeDeclToLinkStr(r, declaration),
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

func TypeDeclToStr(decl *parser.TypeDeclaration) string {
	qname := resolver.NewQualifiedNameFromLocalName(decl.Name)
	tmp := qname.Name()
	if len(decl.Params) > 0 {
		tmp += "<"
		for i, param := range decl.Params {
			tmp += TypeDeclToStr(param)
			if i < len(decl.Params)-1 {
				tmp += ", "
			}
		}
		tmp += ">"
	}

	return tmp
}

func typeDeclToLinkStr(r *resolver.Resolver, decl *parser.TypeDeclaration) string {
	qname := resolver.NewQualifiedNameFromLocalName(decl.Name)
	tmp := qname.Name()
	if len(decl.Params) > 0 {
		tmp += "<"
		for i, param := range decl.Params {
			tmp += typeDeclToLinkStr(r, param)
			if i < len(decl.Params)-1 {
				tmp += ", "
			}
		}
		tmp += ">"
	} else {
		defs := r.Resolve(qname)
		if len(defs) > 0 {
			// TODO this does not work properly in vsc, see also https://github.com/doxygen/doxygen/issues/7421
			tmp = "[[#" + qname.String() + " " + qname.Name() + "]]"
		}

	}

	return tmp
}

func record2Str(data *parser.Struct) string {
	tmp := data.Name.Value + " = \n"
	for i, declaration := range data.Fields {
		tmp += TypeDeclToStr(declaration)
		if i < len(data.Fields)-1 {
			tmp += "\nund "
		}
	}

	return tmp
}

func choice2Str(data *parser.Choice) string {
	tmp := data.Name.Value + " = \n"
	for i, declaration := range data.Choices {
		tmp += TypeDeclToStr(declaration)
		if i < len(data.Choices)-1 {
			tmp += "\noder "
		}
	}

	return tmp
}
