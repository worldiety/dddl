package resolver

import (
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/slices"
	"strings"
)

type Resolver struct {
	ws       *parser.Workspace
	typeDefs map[FullQualifiedName][]*parser.TypeDefinition
	contexts []*Context
}

func NewResolver(ws *parser.Workspace) *Resolver {
	r := &Resolver{typeDefs: map[FullQualifiedName][]*parser.TypeDefinition{}, ws: ws}
	r.initTypeDefLookup()
	r.initContexts()

	return r
}

func (r *Resolver) Workspace() *parser.Workspace {
	return r.ws
}

func (r *Resolver) initTypeDefLookup() {
	parser.MustWalk(r.ws, func(n parser.Node) error {
		if n, ok := n.(*parser.TypeDefinition); ok {
			name := NewQualifiedNameFromNamedType(n.Type)
			defs := r.typeDefs[name]
			defs = append(defs, n)
			r.typeDefs[name] = defs
		}

		return nil
	})
}

type Usage struct {
	Name FullQualifiedName
	Type *parser.TypeDeclaration
}

func (r *Resolver) FindUsages(name FullQualifiedName) []Usage {
	var res []Usage
	for _, definitions := range r.typeDefs {
		for _, definition := range definitions {
			switch t := definition.Type.(type) {
			case *parser.Type:
				if t.Basetype == nil {
					break
				}
				if memberName := NewQualifiedNameFromLocalName(t.Basetype.Name); memberName == name {
					res = append(res, Usage{
						Name: memberName,
						Type: t.Basetype,
					})
				}
			case *parser.Alias:
				if t.BaseType == nil {
					break
				}
				if memberName := NewQualifiedNameFromLocalName(t.BaseType.Name); memberName == name {
					res = append(res, Usage{
						Name: memberName,
						Type: t.BaseType,
					})
				}

			case *parser.Struct:
				for _, field := range t.Fields {
					if memberName := NewQualifiedNameFromLocalName(field.TypeDecl.Name); memberName == name {
						res = append(res, Usage{
							Name: memberName,
							Type: field.TypeDecl,
						})
					}
				}
			case *parser.Choice:
				for _, choice := range t.Choices {
					if memberName := NewQualifiedNameFromLocalName(choice.Choice.Name); memberName == name {
						res = append(res, Usage{
							Name: memberName,
							Type: choice.Choice,
						})
					}
				}
			}
		}
	}

	return res
}

func (r *Resolver) ResolveLocalQualifier(name *parser.QualifiedName) []*parser.TypeDefinition {
	qname := NewQualifiedNameFromLocalName(name)
	return r.typeDefs[qname]
}

func (r *Resolver) Resolve(name FullQualifiedName) []*parser.TypeDefinition {
	return r.typeDefs[name]
}

func (r *Resolver) Guess(nearest parser.Node, name string) []*parser.TypeDefinition {
	// this will try to match exactly
	if strings.Contains(name, ".") {
		defs, ok := r.typeDefs[FullQualifiedName(name)]
		if ok {
			return defs
		}
	}

	// try resolve relatively
	nType := NamedTypeFrom(nearest)
	if nType == nil {
		return nil
	}

	return r.Resolve(NewQualifiedNameFromNamedType(nType).With(name))
}

// Collect returns a sorted slice of the concrete NamedTypes.
func Collect[T parser.NamedType](r *Resolver) []T {
	var res []T
	for _, definitions := range r.typeDefs {
		for _, definition := range definitions {
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
