package parser

type KeywordChoice struct {
	node
	Keyword string `@("enum" | "Auswahl")`
}

type Choice struct {
	node
	KeywordChoice *KeywordChoice     `@@`
	Name          *Name              `@@`
	Choices       []*TypeDeclaration `("{" @@ ((","|"or"|"oder") @@)+ "}" )?`
}

func (n *Choice) GetName() *Name {
	return n.Name
}

func (n *Choice) Children() []Node {
	res := sliceOf(n.KeywordChoice, n.Name)
	for _, choice := range n.Choices {
		res = append(res, choice)
	}
	return res
}

func (*Choice) namedType() {}
