package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
	"html/template"
)

func newTypesFromAliases(parent *Context, r *resolver.Resolver, aliases []*parser.Alias) []*Type {
	var res []*Type
	for _, alias := range aliases {
		res = append(res, newTypeFromAlias(parent, r, alias))
	}

	return res
}

func newTypeFromAlias(parent *Context, r *resolver.Resolver, typ *parser.Alias) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Parent:     parent,
		Category:   "Synonym",
		Name:       typ.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(typ).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ)))
	if err != nil {
		slog.Error("failed to convert alias to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}
