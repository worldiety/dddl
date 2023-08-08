package parser

type KeywordStruct struct {
	node
	Keyword string `@("data" | "Daten")`
}

type Struct struct {
	node
	KeywordStruct *KeywordStruct     `@@`
	Name          *Name              `@@`
	Fields        []*TypeDeclaration `( "{" @@ (("," | "und" ) @@)* "}" )?`
}

func (n *Struct) GetKeyword() string {
	return n.KeywordStruct.Keyword
}

func (n *Struct) Children() []Node {
	res := sliceOf(n.KeywordStruct, n.Name)
	for _, field := range n.Fields {
		res = append(res, field)
	}
	return res
}

func (n *Struct) GetName() *Name {
	return n.Name
}

func (*Struct) namedType() {}
