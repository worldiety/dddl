package parser

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Workspace is nothing to be parsed, but can be used to aggregate multiple documents into something bigger and
// apply the linters on that.
type Workspace struct {
	node
	Documents map[string]*Doc
}

func (n *Workspace) ResolveData(name *Ident) *Data {
	decl := n.ResolveTypeDeclaration(name)
	if data, ok := decl.(*Data); ok {
		return data
	}

	return nil
}

func (n *Workspace) ResolveTypeDeclaration(name *Ident) Declaration {
	ctx := ContextOf(name)

	// lookup in ctx
	if ctx != nil {
		for _, element := range ctx.Elements {
			if element.Name().Value == name.Value {
				if element.DataType != nil {
					return element.DataType
				}

				if element.Workflow != nil {
					return element.Workflow
				}
			}
		}
	}

	// lookup in shared kernel

	for _, doc := range n.Documents {
		for _, definition := range doc.Definitions {
			if definition.TypeDefinition != nil {

				if definition.TypeDefinition.Name().Value == name.Value {
					if definition.TypeDefinition.DataType != nil {
						return definition.TypeDefinition.DataType
					}

					if definition.TypeDefinition.Workflow != nil {
						return definition.TypeDefinition.Workflow
					}
				}
			}
		}
	}

	return nil
}

// Resolve takes the Ident and tries to find the most local declaration:
//   - if Ident is used within a context and context provides a workflow or data declaration, the context qualifier is
//     returned.
//   - if no Ident declaration is found in the context, the anonymous (== shared kernel) space is resolved.
//
// If no declaration was found, false and the expected nearest declaration qualifier is returned.
func (n *Workspace) Resolve(name *Ident) (Qualifier, bool) {
	expectedNearest := Qualifier{Name: name}
	root := name.Parent()
	for root != nil {
		if ctx, ok := root.(*Context); ok {
			nodes := ctx.DeclarationsByName(name.Value)
			if len(nodes) == 0 {
				expectedNearest.Context = ctx
				break
			}

			return Qualifier{
				Context: ctx,
				Name:    name,
			}, true
		}

		root = root.Parent()
	}

	// search "shared kernel" space
	for _, doc := range n.Docs() {
		for _, definition := range doc.Definitions {
			if definition.TypeDefinition != nil {
				if definition.TypeDefinition.Name().Value == name.Value {
					return Qualifier{
						Context: nil,
						Name:    name,
					}, true
				}
			}
		}
	}

	// check if that guy is a context
	for _, doc := range n.Docs() {
		for _, definition := range doc.Definitions {
			if definition.Context != nil {
				if definition.Context.Name.Value == name.Value {
					return Qualifier{
						Context: nil,
						Name:    name,
					}, true
				}
			}
		}
	}

	return expectedNearest, false
}

func (n *Workspace) Contexts() []*Context {
	var tmp []*Context
	for _, n := range n.Children() {
		if c, ok := n.(*Context); ok {
			tmp = append(tmp, c)
		}
	}

	return tmp
}

// Docs returns a stable sorted list of Documents.
func (n *Workspace) Docs() []*Doc {
	var res []*Doc
	keys := maps.Keys(n.Documents)
	slices.Sort(keys)
	for _, key := range keys {
		res = append(res, n.Documents[key])
	}

	return res
}

// UnboundTypeDefinitions are type definitions without an enclosing bounded context.
func (n *Workspace) UnboundTypeDefinitions() []*TypeDecl {
	var tmp []*TypeDecl
	for _, n := range n.Children() {
		if c, ok := n.(*TypeDecl); ok {
			tmp = append(tmp, c)
		}
	}

	return tmp
}

// Children returns a stable sorted list of Documents.
func (n *Workspace) Children() []Node {
	var res []Node

	for _, doc := range n.Docs() {
		res = append(res, doc)
	}

	return res
}
