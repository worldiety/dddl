package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
)

func newTypesFromRecordsInContext(context *Context, r *resolver.Resolver, records []*parser.Struct) []*Type {
	var res []*Type
	for _, record := range records {
		res = append(res, newTypeFromRecordInContext(context, r, record))
	}

	return res
}

func newTypeFromRecordInContext(context *Context, r *resolver.Resolver, record *parser.Struct) *Type {
	typeDef := parser.TypeDefinitionFrom(record)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Context:    context,
		Category:   "Datenverbundtyp",
		Name:       record.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(record).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, record, puml.NewRFlags(record)))
	if err != nil {
		slog.Error("failed to convert data to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)
	data.Usages = newUsages(r, record)

	return data
}

func newTypesFromRecordsInAggregate(aggregate *Aggregate, r *resolver.Resolver, records []*parser.Struct) []*Type {
	var res []*Type
	for _, record := range records {
		res = append(res, newTypeFromRecordsInAggregate(aggregate, r, record))
	}

	return res
}

func newTypeFromRecordsInAggregate(aggregate *Aggregate, r *resolver.Resolver, record *parser.Struct) *Type {
	typeDef := parser.TypeDefinitionFrom(record)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Aggregate:  aggregate,
		Category:   "Datenverbundtyp",
		Name:       record.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(record).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, record, puml.NewRFlags(record)))
	if err != nil {
		slog.Error("failed to convert data to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}
