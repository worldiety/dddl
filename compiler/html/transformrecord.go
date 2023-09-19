package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromRecords(parent any, r *resolver.Resolver, model PreviewModel, records []*parser.Struct) []*Type {
	var res []*Type
	for _, record := range records {
		res = append(res, newTypeFromRecord(parent, r, model, record))
	}

	return res
}

func newTypeFromRecord(parent any, r *resolver.Resolver, model PreviewModel, record *parser.Struct) *Type {
	typeDef := parser.TypeDefinitionFrom(record)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Type{
		Node:                record,
		Parent:              parent,
		Category:            "Datenverbundtyp",
		Name:                record.Name.Value,
		Ref:                 resolver.NewQualifiedNameFromNamedType(record).String(),
		Definition:          def,
		SVG:                 "",
		WorkPackageName:     parser.FindAnnotation[*parser.WorkPackageAnnotation](record).GetName(),
		WorkPackageRequires: parser.FindAnnotation[*parser.WorkPackageAnnotation](record).GetRequires(),
		WorkPackageDuration: parser.FindAnnotation[*parser.WorkPackageAnnotation](record).GetDuration(),
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, record, puml.NewRFlags(record)))
	if err != nil {
		slog.Error("failed to convert data to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)
	data.Usages = newUsages(r, record)

	return data
}
