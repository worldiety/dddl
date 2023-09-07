package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromTypes(parent any, r *resolver.Resolver, types []*parser.Type) []*Type {
	var res []*Type
	for _, typ := range types {
		res = append(res, newTypeFromType(parent, r, typ))
	}

	return res
}

func newTypeFromType(parent any, r *resolver.Resolver, typ *parser.Type) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Node:       typ,
		Parent:     parent,
		Category:   "Basistyp",
		Name:       typ.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(typ).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ)))
	if err != nil {
		slog.Error("failed to convert type to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)
	data.Usages = newUsages(r, typ)

	return data
}
