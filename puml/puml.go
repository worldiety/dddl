package puml

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
	"math"
)

type RFlags struct {
	MainType parser.NamedType
	Visited  map[parser.NamedType]bool
	Depth    int
}

func NewRFlags(mainType parser.NamedType) RFlags {
	return RFlags{
		MainType: mainType,
		Visited:  map[parser.NamedType]bool{},
		Depth:    2,
	}
}

func (r RFlags) WithMaxDepth() RFlags {
	r.Depth = math.MaxInt
	return r
}

func RenderNamedType(r *resolver.Resolver, namedType parser.NamedType, flags RFlags) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	diag.BackgroundColor = "#00000000"
	if _, ok := flags.Visited[namedType]; ok {
		return diag
	}

	if flags.Depth <= 0 {
		return diag
	}

	flags.Depth--

	flags.Visited[namedType] = false

	switch t := namedType.(type) {
	case *parser.Struct:
		diag.Add(Record(r, t, flags).Renderables...)
	case *parser.Choice:
		diag.Add(Choice(r, t, flags).Renderables...)
		if len(t.Choices) > 5 {
			diag.SkinParams = append(diag.SkinParams, "left to right direction")
		}
	case *parser.Type:
		diag.Add(Type(r, t, flags).Renderables...)
	case *parser.Alias:
		diag.Add(Alias(r, t, flags).Renderables...)
	case *parser.Function:
		tmp := Func(r, t, flags)
		flags.Visited[namedType] = true
		diag.DefaultTextAlignment = tmp.DefaultTextAlignment
		diag.Add(tmp.Renderables...)
	}

	return diag
}
