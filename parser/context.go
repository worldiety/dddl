package parser

type Context struct {
	node
	KeywordContext *KeywordContext   `@@`
	Name           *Name             `@@`
	Definitions    []*TypeDefinition `("{"@@*"}")?`
}

func (n *Context) GetName() *Name {
	return n.Name
}

func (*Context) namedType() {}

func ContextOf(root Node) *Context {
	for root != nil {
		if wf, ok := root.(*Context); ok {
			return wf
		}
		root = root.Parent()
	}

	return nil
}

func (n *Context) Children() []Node {
	res := sliceOf(n.KeywordContext, n.Name)
	for _, definition := range n.Definitions {
		res = append(res, definition)
	}
	return res
}
