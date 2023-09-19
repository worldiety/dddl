package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromFuncs(parent any, r *resolver.Resolver, model PreviewModel, funcs []*parser.Function) []*Type {
	var res []*Type
	for _, f := range funcs {
		res = append(res, newTypeFromFunc(parent, r, model, f))
	}

	return res
}

func newTypeFromFunc(parent any, r *resolver.Resolver, model PreviewModel, typ *parser.Function) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Type{
		Node:                typ,
		Parent:              parent,
		Category:            typ.KeywordFn.Keyword,
		Name:                typ.Name.Value,
		Ref:                 resolver.NewQualifiedNameFromNamedType(typ).String(),
		Definition:          def,
		SVG:                 "",
		WorkPackageName:     parser.FindAnnotation[*parser.WorkPackageAnnotation](typ).GetName(),
		WorkPackageRequires: parser.FindAnnotation[*parser.WorkPackageAnnotation](typ).GetRequires(),
		WorkPackageDuration: parser.FindAnnotation[*parser.WorkPackageAnnotation](typ).GetDuration(),
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ)))
	if err != nil {
		slog.Error("failed to convert func to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}
