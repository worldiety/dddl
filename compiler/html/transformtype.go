package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromTypes(ctx *plantuml.PreflightContext, parent any, r *resolver.Resolver, model PreviewModel, types []*parser.Type) []*Type {
	var res []*Type
	for _, typ := range types {
		res = append(res, newTypeFromType(ctx, parent, r, model, typ))
	}

	return res
}

func newTypeFromType(ctx *plantuml.PreflightContext, parent any, r *resolver.Resolver, model PreviewModel, typ *parser.Type) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Type{
		Node:                typ,
		Parent:              parent,
		Category:            "Basistyp",
		Name:                typ.Name.Value,
		Ref:                 resolver.NewQualifiedNameFromNamedType(typ).String(),
		Definition:          def,
		SVG:                 "",
		WorkPackageName:     parser.FindAnnotation[*parser.WorkPackageAnnotation](typ).GetName(),
		WorkPackageRequires: parser.FindAnnotation[*parser.WorkPackageAnnotation](typ).GetRequires(),
		WorkPackageDuration: parser.FindAnnotation[*parser.WorkPackageAnnotation](typ).GetDuration(),
	}

	svg, err := plantuml.RenderLocalWithPreflight(ctx, "svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ)))
	if err != nil {
		slog.Error("failed to convert type to puml", slog.Any("err", err))
	}

	svgX, err := plantuml.RenderLocalWithPreflight(ctx, "svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ).WithMaxDepth()))
	if err != nil {
		slog.Error("failed to convert alias to puml", slog.Any("err", err))
	}

	data.SVGExtended = template.HTML(svgX)
	data.SVG = template.HTML(svg)
	if data.SVGExtended == data.SVG {
		data.SVGExtended = ""
	}

	data.Usages = newUsages(r, typ)

	return data
}
