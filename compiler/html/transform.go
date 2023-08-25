package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

func transform(rslv *resolver.Resolver) *Doc {

	doc := &Doc{}
	for _, rCtx := range rslv.Contexts() {
		ctx := &Context{Name: rCtx.Name}
		ctx.ShortDef = markdown(rCtx.ShortString())
		ctx.Ref = resolver.NewQualifiedNameFromNamedType(rCtx.Fragments[0]).String()
		ctx.Aggregates = append(ctx.Aggregates, newTypesFromAggregateInContext(ctx, rslv, resolver.CollectFromContext[*parser.Aggregate](rCtx), rCtx)...)
		ctx.Types = append(ctx.Types, newTypesFromRecordsInContext(ctx, rslv, resolver.CollectFromContext[*parser.Struct](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromChoiceInContext(ctx, rslv, resolver.CollectFromContext[*parser.Choice](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromTypesInContext(ctx, rslv, resolver.CollectFromContext[*parser.Type](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromAliasInContext(ctx, rslv, resolver.CollectFromContext[*parser.Alias](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromFuncsInContext(ctx, rslv, resolver.CollectFromContext[*parser.Function](rCtx))...)

		ctx.Definition = markdown(rCtx.Description)

		doc.Contexts = append(doc.Contexts, ctx)

	}

	return doc
}
