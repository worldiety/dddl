package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromAlias(ctx *plantuml.PreflightContext, parent any, r *resolver.Resolver, model PreviewModel, aliases []*parser.Alias) []*Type {
	var res []*Type
	for _, alias := range aliases {
		res = append(res, newTypeFromAlias(ctx, parent, r, model, alias))
	}

	return res
}

func newTypeFromAlias(ctx *plantuml.PreflightContext, parent any, r *resolver.Resolver, model PreviewModel, typ *parser.Alias) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Type{
		Node:                typ,
		Parent:              parent,
		Category:            "Synonym",
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
		slog.Error("failed to convert alias to puml", slog.Any("err", err))
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
