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
	Error     error
}

func (n *Workspace) ResolveData(name *Ident) *Data {
	decl := n.ResolveTypeDeclaration(name)
	if data, ok := decl.(*Data); ok {
		return data
	}

	return nil
}

func (n *Workspace) ContextByName(name string) []*Context {
	for _, c := range n.CollectContextChildren() {
		if c.Name == name {
			return c.Contexts
		}
	}

	return nil
}

func (n *Workspace) ResolveTypeDeclaration(name *Ident) Declaration {
	ctx := ContextOf(name)

	// lookup in all equally named ctx definitions
	if ctx != nil {
		for _, ctx := range n.ContextByName(name.Value) {
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
			for _, ctx := range n.ContextByName(ctx.Name.Value) {
				nodes := ctx.DeclarationsByName(name.Value)
				if len(nodes) == 0 {
					expectedNearest.Context = ctx
					continue
				}

				return Qualifier{
					Context: ctx,
					Name:    name,
				}, true
			}

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

func (n *Workspace) CollectFreeDataOrWorkflow() []DataOrWorkflow {
	var res []DataOrWorkflow
	for _, doc := range n.Documents {
		for _, definition := range doc.Definitions {
			if definition.TypeDefinition != nil {
				if definition.TypeDefinition.DataType != nil {
					res = append(res, definition.TypeDefinition.DataType)
				}

				if definition.TypeDefinition.Workflow != nil {
					res = append(res, definition.TypeDefinition.Workflow)
				}
			}

		}
	}

	return res
}

type ContextCollection[T any] struct {
	Name     string // name of context
	Children []T
	Contexts []*Context
}

func (n *Workspace) CollectContextChildren() []ContextCollection[Node] {
	return CollectContextChildren[Node](n, func(context *Context) []Node {
		var res []Node
		for _, element := range context.Elements {
			if element.Workflow != nil {
				res = append(res, element.Workflow)
			}

			if element.DataType != nil {
				res = append(res, element.DataType)
			}
		}
		return res
	})
}

func CollectContextChildren[T any](ws *Workspace, collect func(ctx *Context) []T) []ContextCollection[T] {
	tmp := map[string][]T{}
	tmp2 := map[string][]*Context{}

	for _, doc := range ws.Documents {
		for _, context := range doc.Contexts() {
			tmp2[context.Name.Value] = append(tmp2[context.Name.Value], context)
			for _, t := range collect(context) {
				tmp[context.Name.Value] = append(tmp[context.Name.Value], t)
			}

		}
	}

	keys := maps.Keys(tmp)
	slices.Sort(keys)

	var res []ContextCollection[T]
	for _, key := range keys {
		collection := ContextCollection[T]{Name: key}
		for _, data := range tmp[key] {
			collection.Children = append(collection.Children, data)
		}
		collection.Contexts = tmp2[key]
		res = append(res, collection)
	}

	return res
}

func (n *Workspace) Contexts() []*Context {
	var res []*Context
	for _, doc := range n.Docs() {
		res = append(res, doc.Contexts()...)
	}

	return res
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
