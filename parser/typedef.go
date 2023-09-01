package parser

type KeyValue struct {
	node
	Key   *Name `@@ ( "="`
	Value *Name `@@ )?`
}

func (n *KeyValue) Children() []Node {
	return sliceOf(n.Key, n.Value)
}

type Annotation struct {
	node
	Name      *Name       `@@`
	KeyValues []*KeyValue `( "("  @@ ("," @@)*  ")" )?`
}

func (n *Annotation) Children() []Node {
	res := sliceOf(n.Name)
	for _, value := range n.KeyValues {
		res = append(res, value)
	}

	return res
}

type TypeDefinition struct {
	node
	// Description may be nil
	Description *Literal      `@@?`
	Annotations []*Annotation `("@" @@)*`
	Type        NamedType     `@@`
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
	for _, tag := range n.Annotations {
		res = append(res, tag)
	}

	res = append(res, n.Type)
	return res
}

type NamedType interface {
	Node
	namedType()
	GetName() *Name
	GetKeyword() string
}
