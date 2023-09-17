package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromAlias(parent any, r *resolver.Resolver, model PreviewModel, aliases []*parser.Alias) []*Type {
	var res []*Type
	for _, alias := range aliases {
		res = append(res, newTypeFromAlias(parent, r, model, alias))
	}

	return res
}

func newTypeFromAlias(parent any, r *resolver.Resolver, model PreviewModel, typ *parser.Alias) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Type{
		Node:       typ,
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
	data.Usages = newUsages(r, typ)

	return data
}

func newUsages(r *resolver.Resolver, namedType parser.NamedType) []Usage {
	var res []Usage
	for _, usage := range r.FindUsages(resolver.NewQualifiedNameFromNamedType(namedType)) {
		res = append(res, Usage{
			Name: usage.Name.Name(),
			Ref:  usage.Name.String(),
		})
	}

	return res
}
