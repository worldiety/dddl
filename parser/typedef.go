package parser

type TypeDefinition struct {
	node
	// Description may be nil
	Description *Literal  `@@?`
	Tags        []*Name   `("@" @@)*`
	Type        NamedType `@@`
}

func TypeDefinitionFrom(n Node) *TypeDefinition {
	for n != nil {
		if td, ok := n.(*TypeDefinition); ok {
			return td
		}

		n = n.Parent()
	}

	return nil
}

func (n *TypeDefinition) Children() []Node {
	res := sliceOf(n.Description)
	for _, tag := range n.Tags {
		res = append(res, tag)
	}

	res = append(res, n.Type)
	return res
}

type NamedType interface {
	Node
	namedType()
	GetName() *Name
}
