package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

func transform(rslv *resolver.Resolver, model PreviewModel) *Doc {

	doc := &Doc{}
	for _, rCtx := range rslv.Contexts() {
		ctx := &Context{Name: rCtx.Name}
		ctx.ShortDef = markdown(rCtx.ShortString(), model)
		ctx.Ref = resolver.NewQualifiedNameFromNamedType(rCtx.Fragments[0]).String()
		ctx.Aggregates = append(ctx.Aggregates, newTypesFromAggregate(ctx, rslv, model, resolver.CollectFromContext[*parser.Aggregate](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromRecords(ctx, rslv, model, resolver.CollectFromContext[*parser.Struct](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromChoice(ctx, rslv, model, resolver.CollectFromContext[*parser.Choice](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromTypes(ctx, rslv, model, resolver.CollectFromContext[*parser.Type](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromAlias(ctx, rslv, model, resolver.CollectFromContext[*parser.Alias](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromFuncs(ctx, rslv, model, resolver.CollectFromContext[*parser.Function](rCtx))...)

		postCategorizeByAnnotations(ctx.Types)
		ctx.Definition = markdown(rCtx.Description, model)

		ctx.WorkPackageName = parser.FindAnnotation[*parser.WorkPackageAnnotation](rCtx.Fragments[0]).GetName()
		ctx.WorkPackageRequires = parser.FindAnnotation[*parser.WorkPackageAnnotation](rCtx.Fragments[0]).GetRequires()
		ctx.WorkPackageDuration = parser.FindAnnotation[*parser.WorkPackageAnnotation](rCtx.Fragments[0]).GetDuration()
		doc.Contexts = append(doc.Contexts, ctx)

	}

	return doc
}

func postCategorizeByAnnotations(types []*Type) {
	// post-process annotations
	for i := range types {
		typeDef := types[i].Node.Parent().(*parser.TypeDefinition)
		evtA := parser.FindAnnotation[*parser.EventAnnotation](typeDef)
		if evtA != nil {
			types[i].Category = "Ereignis"
		}

		errA := parser.FindAnnotation[*parser.ErrorAnnotation](typeDef)
		if errA != nil {
			types[i].Category = "Fehler"
		}

		extA := parser.FindAnnotation[*parser.ExternalSystemAnnotation](typeDef)
		if extA != nil {
			types[i].Category = "Fremdsystem"
		}

		roleA := parser.FindAnnotation[*parser.RoleAnnotation](typeDef)
		if roleA != nil {
			types[i].Category = "Rolle"
		}
	}
}
