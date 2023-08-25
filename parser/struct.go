package parser

type KeywordStruct struct {
	node
	Keyword string `@("data" | "Daten")`
}

type Field struct {
	node
	TypeDecl *TypeDeclaration `@@`
	Alias    *Name            `(("als" | "as") @@)?`
}

// Name returns the fields name, which is either the last segment of the TypeDeclaration name or the optional field
// Alias.
func (n *Field) Name() string {
	if n.Alias != nil {
		return n.Alias.Value
	}

	return n.TypeDecl.Name.Name()
}

func (n *Field) Children() []Node {
	return sliceOf(n.TypeDecl)
}

type Struct struct {
	node
	KeywordStruct *KeywordStruct `@@`
	Name          *Name          `@@`
	Fields        []*Field       `( "{" @@ (("," | "und" ) @@)* "}" )?`
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
