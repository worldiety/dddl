package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromChoiceInContext(context *Context, r *resolver.Resolver, choices []*parser.Choice) []*Type {
	var res []*Type
	for _, choice := range choices {
		res = append(res, newTypeFromChoiceInContext(context, r, choice))
	}

	return res
}

func newTypeFromChoiceInContext(context *Context, r *resolver.Resolver, choice *parser.Choice) *Type {
	typeDef := parser.TypeDefinitionFrom(choice)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Context:    context,
		Category:   "Auswahltyp",
		Name:       choice.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(choice).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, choice, puml.NewRFlags(choice)))
	if err != nil {
		slog.Error("failed to convert choice to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}

func newTypesFromChoiceInAggregate(aggregate *Aggregate, r *resolver.Resolver, choices []*parser.Choice) []*Type {
	var res []*Type
	for _, choice := range choices {
		res = append(res, newTypeFromChoiceInAggregate(aggregate, r, choice))
	}

	return res
}

func newTypeFromChoiceInAggregate(aggregate *Aggregate, r *resolver.Resolver, choice *parser.Choice) *Type {
	typeDef := parser.TypeDefinitionFrom(choice)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Aggregate:  aggregate,
		Category:   "Auswahltyp",
		Name:       choice.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(choice).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, choice, puml.NewRFlags(choice)))
	if err != nil {
		slog.Error("failed to convert choice to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}
