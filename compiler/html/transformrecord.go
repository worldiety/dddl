package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
	"html/template"
)

func newTypesFromRecords(parent *Context, r *resolver.Resolver, records []*parser.Struct) []*Type {
	var res []*Type
	for _, record := range records {
		res = append(res, newTypeFromRecord(parent, r, record))
	}

	return res
}

func newTypeFromRecord(parent *Context, r *resolver.Resolver, record *parser.Struct) *Type {
	typeDef := parser.TypeDefinitionFrom(record)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Parent:     parent,
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
