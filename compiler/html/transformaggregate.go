package html

import (
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

func newTypesFromAggregateInContext(context *Context, r *resolver.Resolver, aggregates []*parser.Aggregate, rCtx *resolver.Context) []*Aggregate {
	var res []*Aggregate
	for _, aggregate := range aggregates {
		res = append(res, newTypeFromAggregateInContext(context, r, aggregate, rCtx))
	}

	return res
}

func newTypeFromAggregateInContext(context *Context, r *resolver.Resolver, aggregate *parser.Aggregate, rCtx *resolver.Context) *Aggregate {
	typeDef := parser.TypeDefinitionFrom(aggregate)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Aggregate{
		Context:    context,
		Category:   "Aggregattyp",
		Name:       aggregate.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(aggregate).String(),
		Definition: def,
	}

	data.Types = append(data.Types, newTypesFromRecordsInAggregate(data, r, resolver.CollectFromAggregate[*parser.Struct](aggregate))...)
	data.Types = append(data.Types, newTypesFromChoiceInAggregate(data, r, resolver.CollectFromAggregate[*parser.Choice](aggregate))...)
	data.Types = append(data.Types, newTypesFromTypesInAggregate(data, r, resolver.CollectFromAggregate[*parser.Type](aggregate))...)
	data.Types = append(data.Types, newTypesFromAliasInAggregate(data, r, resolver.CollectFromAggregate[*parser.Alias](aggregate))...)
	data.Types = append(data.Types, newTypesFromFuncsInAggregate(data, r, resolver.CollectFromAggregate[*parser.Function](aggregate))...)

	return data
}
