package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

func transform(pWS *parser.Workspace) *Doc {
	rslv := resolver.NewResolver(pWS)
	doc := &Doc{}
	for _, rCtx := range rslv.Contexts() {
		ctx := &Context{Name: rCtx.Name}
		ctx.ShortDef = markdown(rCtx.ShortString())
		ctx.Ref = resolver.NewQualifiedNameFromNamedType(rCtx.Fragments[0]).String()
		ctx.Types = append(ctx.Types, newTypesFromRecords(ctx, rslv, resolver.CollectFromContext[*parser.Struct](rCtx))...)
		ctx.Types = append(ctx.Types, newTypesFromChoice(ctx, rslv, resolver.CollectFromContext[*parser.Choice](rCtx))...)
		//	ctx.Types = append(ctx.Types, newTypesFromRecords(ctx, resolver.CollectFromContext[*parser.Choice](rCtx))...)
		//	ctx.Types = append(ctx.Types, newTypesFromRecords(ctx, resolver.CollectFromContext[*parser.Function](rCtx))...)

		ctx.Definition = markdown(rCtx.Description)

		doc.Contexts = append(doc.Contexts, ctx)

	}

	return doc
}
