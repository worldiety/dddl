package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromChoice(ctx *plantuml.PreflightContext, parent any, r *resolver.Resolver, model PreviewModel, choices []*parser.Choice) []*Type {
	var res []*Type
	for _, choice := range choices {
		res = append(res, newTypeFromChoice(ctx, parent, r, model, choice))
	}

	return res
}

func newTypeFromChoice(ctx *plantuml.PreflightContext, parent any, r *resolver.Resolver, model PreviewModel, choice *parser.Choice) *Type {
	typeDef := parser.TypeDefinitionFrom(choice)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Type{
		Node:                choice,
		Parent:              parent,
		Category:            "Auswahltyp",
		Name:                choice.Name.Value,
		Ref:                 resolver.NewQualifiedNameFromNamedType(choice).String(),
		Definition:          def,
		SVG:                 "",
		WorkPackageName:     parser.FindAnnotation[*parser.WorkPackageAnnotation](choice).GetName(),
		WorkPackageRequires: parser.FindAnnotation[*parser.WorkPackageAnnotation](choice).GetRequires(),
		WorkPackageDuration: parser.FindAnnotation[*parser.WorkPackageAnnotation](choice).GetDuration(),
	}

	svg, err := plantuml.RenderLocalWithPreflight(ctx, "svg", puml.RenderNamedType(r, choice, puml.NewRFlags(choice)))
	if err != nil {
		slog.Error("failed to convert choice to puml", slog.Any("err", err))
	}

	svgX, err := plantuml.RenderLocalWithPreflight(ctx, "svg", puml.RenderNamedType(r, choice, puml.NewRFlags(choice).WithMaxDepth()))
	if err != nil {
		slog.Error("failed to convert alias to puml", slog.Any("err", err))
	}

	data.SVGExtended = template.HTML(svgX)
	data.SVG = template.HTML(svg)
	if data.SVGExtended == data.SVG {
		data.SVGExtended = ""
	}

	data.Usages = newUsages(r, choice)

	return data
}
