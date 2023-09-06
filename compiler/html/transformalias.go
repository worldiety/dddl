package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromAliasInContext(context *Context, r *resolver.Resolver, aliases []*parser.Alias) []*Type {
	var res []*Type
	for _, alias := range aliases {
		res = append(res, newTypeFromAliasInContext(context, r, alias))
	}

	return res
}

func newTypeFromAliasInContext(context *Context, r *resolver.Resolver, typ *parser.Alias) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Context:    context,
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

func newTypesFromAliasInAggregate(aggregate *Aggregate, r *resolver.Resolver, aliases []*parser.Alias) []*Type {
	var res []*Type
	for _, alias := range aliases {
		res = append(res, newTypeFromAliasInAggregate(aggregate, r, alias))
	}

	return res
}

func newTypeFromAliasInAggregate(aggregate *Aggregate, r *resolver.Resolver, typ *parser.Alias) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Aggregate:  aggregate,
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
