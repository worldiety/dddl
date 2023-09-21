package html

import (
	"github.com/worldiety/dddl/plantuml"
	"html/template"

	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

func newTypesFromAggregate(ctx *plantuml.PreflightContext, context *Context, r *resolver.Resolver, model PreviewModel, aggregates []*parser.Aggregate) []*Aggregate {
	var res []*Aggregate
	for _, aggregate := range aggregates {
		a := newTypeFromAggregate(ctx, context, r, model, aggregate)
		postCategorizeByAnnotations(a.Types)
		res = append(res, a)
	}

	return res
}

func newTypeFromAggregate(ctx *plantuml.PreflightContext, context *Context, r *resolver.Resolver, model PreviewModel, aggregate *parser.Aggregate) *Aggregate {
	typeDef := parser.TypeDefinitionFrom(aggregate)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value, model)
	}

	data := &Aggregate{
		Context:             context,
		Category:            "Aggregattyp",
		Name:                aggregate.Name.Value,
		Ref:                 resolver.NewQualifiedNameFromNamedType(aggregate).String(),
		Definition:          def,
		WorkPackageName:     parser.FindAnnotation[*parser.WorkPackageAnnotation](aggregate).GetName(),
		WorkPackageRequires: parser.FindAnnotation[*parser.WorkPackageAnnotation](aggregate).GetRequires(),
		WorkPackageDuration: parser.FindAnnotation[*parser.WorkPackageAnnotation](aggregate).GetDuration(),
	}

	data.Types = append(data.Types, newTypesFromRecords(ctx, data, r, model, resolver.CollectFromAggregate[*parser.Struct](aggregate))...)
	data.Types = append(data.Types, newTypesFromChoice(ctx, data, r, model, resolver.CollectFromAggregate[*parser.Choice](aggregate))...)
	data.Types = append(data.Types, newTypesFromTypes(ctx, data, r, model, resolver.CollectFromAggregate[*parser.Type](aggregate))...)
	data.Types = append(data.Types, newTypesFromAlias(ctx, data, r, model, resolver.CollectFromAggregate[*parser.Alias](aggregate))...)
	data.Types = append(data.Types, newTypesFromFuncs(ctx, data, r, model, resolver.CollectFromAggregate[*parser.Function](aggregate))...)

	return data
}
