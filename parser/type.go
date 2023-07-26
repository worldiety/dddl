package parser

type KeywordType struct {
	node
	Keyword string `@("type" | "Typ")`
}

type Type struct {
	node
	KeywordType *KeywordType     `@@`
	Name        *Name            `@@`
	Basetype    *TypeDeclaration `( "=" @@ )?`
}

func (n *Type) GetName() *Name {
	return n.Name
}

func (*Type) namedType() {}
