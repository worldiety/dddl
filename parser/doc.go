package parser

type Doc struct {
	node
	Types []*TypeDefinition `@@*`
}

func (n *Doc) NewVirtualContext(name string) *Context {
	return &Context{
		node: n.node,
		KeywordContext: &KeywordContext{
			node:    n.node,
			Keyword: "context",
		},
		Name: &Name{
			node:  n.node,
			Value: name,
		},
	}
}

func (n *Doc) Children() []Node {
	var res []Node
	for _, t := range n.Types {
		res = append(res, t)
	}

	return res
}
