package resolver

import (
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"strings"
)

type Context struct {
	Name        string
	Description string
	Fragments   []*parser.Context
}

func (c *Context) ShortString() string {
	if len(c.Description) < 200 {
		return c.Description
	}

	if idx := strings.Index(c.Description, "."); idx > 0 {
		return c.Description[:idx+1]
	}

	return c.Description
}

func (r *Resolver) Contexts() []*Context {
	return r.contexts
}

// Context returns the resolved and aggregated Context or nil.
func (r *Resolver) Context(name string) *Context {
	for _, context := range r.contexts {
		if context.Name == name {
			return context
		}
	}

	return nil
}

func (r *Resolver) initContexts() {
	tmp := map[string][]*parser.Context{}
	parser.MustWalk(r.ws, func(n parser.Node) error {
		if ctx, ok := n.(*parser.Context); ok {
			tmp[ctx.Name.Value] = append(tmp[ctx.Name.Value], ctx)
		}
		return nil
	})

	keys := maps.Keys(tmp)
	slices.Sort(keys)
	for _, key := range keys {
		context := &Context{Name: key}
		for _, ctx := range tmp[key] {
			context.Fragments = append(context.Fragments, ctx)
		}

		slices.SortFunc(context.Fragments, func(a, b *parser.Context) bool {
			return a.Pos.Filename < b.Pos.Filename
		})

		desc := ""
		for _, fragment := range context.Fragments {
			if typeDef, ok := fragment.Parent().(*parser.TypeDefinition); ok {
				if typeDef.Description != nil {
					desc += parser.TextOf(typeDef.Description.Value)
					desc += "\n\n"
				}

			}
		}

		desc = strings.TrimSpace(desc)
		context.Description = desc

		r.contexts = append(r.contexts, context)
	}
}

// CollectFromContext returns a sorted slice of the concrete NamedTypes.
func CollectFromContext[T parser.NamedType](ctx *Context) []T {
	var res []T
	for _, fragment := range ctx.Fragments {
		for _, definition := range fragment.Definitions {
			if t, ok := definition.Type.(T); ok {
				res = append(res, t)
			}
		}
	}

	slices.SortFunc(res, func(a, b T) bool {
		return a.GetName().Value < b.GetName().Value
	})

	return res
}
