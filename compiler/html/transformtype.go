package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
	"html/template"
)

func newTypesFromTypes(parent *Context, r *resolver.Resolver, types []*parser.Type) []*Type {
	var res []*Type
	for _, typ := range types {
		res = append(res, newTypeFromType(parent, r, typ))
	}

	return res
}

func newTypeFromType(parent *Context, r *resolver.Resolver, typ *parser.Type) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
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

	return data
}
