package parser

type Doc struct {
	node
	Types []*TypeDefinition `@@*`
}

func (n *Doc) Children() []Node {
	var res []Node
	for _, t := range n.Types {
		res = append(res, t)
	}

	return res
}
