package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromFuncsInContext(context *Context, r *resolver.Resolver, funcs []*parser.Function) []*Type {
	var res []*Type
	for _, f := range funcs {
		res = append(res, newTypeFromFuncInContext(context, r, f))
	}

	return res
}

func newTypeFromFuncInContext(context *Context, r *resolver.Resolver, typ *parser.Function) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Context:    context,
		Category:   typ.KeywordFn.Keyword,
		Name:       typ.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(typ).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ)))
	if err != nil {
		slog.Error("failed to convert func to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}

func newTypesFromFuncsInAggregate(aggregate *Aggregate, r *resolver.Resolver, funcs []*parser.Function) []*Type {
	var res []*Type
	for _, f := range funcs {
		res = append(res, newTypeFromFuncInAggregate(aggregate, r, f))
	}

	return res
}

func newTypeFromFuncInAggregate(aggregate *Aggregate, r *resolver.Resolver, typ *parser.Function) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Aggregate:  aggregate,
		Category:   typ.KeywordFn.Keyword,
		Name:       typ.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(typ).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ)))
	if err != nil {
		slog.Error("failed to convert func to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}
