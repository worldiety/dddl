package parser

type KeywordAlias struct {
	node
	Keyword string `@("alias" | "Synonym")`
}

type Alias struct {
	node
	KeywordAlias *KeywordAlias    `@@`
	Name         *Name            `@@`
	BaseType     *TypeDeclaration `( "=" @@ )?`
}

func (n *Alias) GetName() *Name {
	return n.Name
}

func (n *Alias) Children() []Node {
	return sliceOf(n.KeywordAlias, n.Name, n.BaseType)
}

func (*Alias) namedType() {}
