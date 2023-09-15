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
		ctx.Aggregates = append(ctx.Aggregates, newTypesFromAggregate(ctx, rslv, resolver.CollectFromContext[*parser.Aggregate](rCtx), rCtx)...)
		ctx.Types = append(ctx.Types, newTypesFromRecords(ctx, rslv, resolver.CollectFromContext[*parser.Struct](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromChoice(ctx, rslv, resolver.CollectFromContext[*parser.Choice](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromTypes(ctx, rslv, resolver.CollectFromContext[*parser.Type](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromAlias(ctx, rslv, resolver.CollectFromContext[*parser.Alias](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromFuncs(ctx, rslv, resolver.CollectFromContext[*parser.Function](rCtx))...)

		postCategorizeByAnnotations(ctx.Types)
		ctx.Definition = markdown(rCtx.Description)

		doc.Contexts = append(doc.Contexts, ctx)

	}

	return doc
}

func postCategorizeByAnnotations(types []*Type) {
	// post-process annotations
	for i := range types {
		typeDef := types[i].Node.Parent().(*parser.TypeDefinition)
		evtA, _ := parser.ParseEventAnnotation(typeDef)
		if evtA != nil {
			types[i].Category = "Ereignis"
		}

		errA, _ := parser.ParseErrorAnnotation(typeDef)
		if errA != nil {
			types[i].Category = "Fehler"
		}

		extA, _ := parser.ParseExternalSystemAnnotation(typeDef)
		if extA != nil {
			types[i].Category = "Fremdsystem"
		}
	}
}
