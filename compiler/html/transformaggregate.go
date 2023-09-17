package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

func newTypesFromAggregate(context *Context, r *resolver.Resolver, model PreviewModel, aggregates []*parser.Aggregate) []*Aggregate {
	var res []*Aggregate
	for _, aggregate := range aggregates {
		a := newTypeFromAggregate(context, r, model, aggregate)
		postCategorizeByAnnotations(a.Types)
		res = append(res, a)
	}

	return res
}

func newTypeFromAggregate(context *Context, r *resolver.Resolver, model PreviewModel, aggregate *parser.Aggregate) *Aggregate {
	typeDef := parser.TypeDefinitionFrom(aggregate)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Aggregate{
		Context:    context,
		Category:   "Aggregattyp",
		Name:       aggregate.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(aggregate).String(),
		Definition: def,
	}

	data.Types = append(data.Types, newTypesFromRecords(data, r, model, resolver.CollectFromAggregate[*parser.Struct](aggregate))...)
	data.Types = append(data.Types, newTypesFromChoice(data, r, model, resolver.CollectFromAggregate[*parser.Choice](aggregate))...)
	data.Types = append(data.Types, newTypesFromTypes(data, r, model, resolver.CollectFromAggregate[*parser.Type](aggregate))...)
	data.Types = append(data.Types, newTypesFromAlias(data, r, model, resolver.CollectFromAggregate[*parser.Alias](aggregate))...)
	data.Types = append(data.Types, newTypesFromFuncs(data, r, model, resolver.CollectFromAggregate[*parser.Function](aggregate))...)

	return data
}
