package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/slog"
	"html/template"
)

func newTypesFromFuncs(parent *Context, r *resolver.Resolver, funcs []*parser.Function) []*Type {
	var res []*Type
	for _, f := range funcs {
		res = append(res, newTypeFromFunc(parent, r, f))
	}

	return res
}

func newTypeFromFunc(parent *Context, r *resolver.Resolver, typ *parser.Function) *Type {
	typeDef := parser.TypeDefinitionFrom(typ)
	var def template.HTML
	if typeDef.Description != nil {
		def = markdown(typeDef.Description.Value)
	}

	data := &Type{
		Parent:     parent,
		Category:   typ.KeywordFn.Keyword,
		Name:       typ.Name.Value,
		Ref:        resolver.NewQualifiedNameFromNamedType(typ).String(),
		Definition: def,
		SVG:        "",
	}

	svg, err := plantuml.RenderLocal("svg", puml.RenderNamedType(r, typ, puml.NewRFlags(typ)))
	if err != nil {
		slog.Error("failed to convert func to puml", slog.Any("err", err))
	}

	data.SVG = template.HTML(svg)

	return data
}
