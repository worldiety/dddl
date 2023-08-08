package parser

type Type struct {
	node
	KeywordType *KeywordType     `@@`
	Name        *Name            `@@`
	Basetype    *TypeDeclaration `( "=" @@ )?`
}

func (n *Type) Children() []Node {
	return sliceOf(n.KeywordType, n.Name, n.Basetype)
}

func (n *Type) GetKeyword() string {
	return n.KeywordType.Keyword
}

func (n *Type) GetName() *Name {
	return n.Name
}

func (*Type) namedType() {}
