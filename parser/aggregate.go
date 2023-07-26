package parser

type KeywordAggregate struct {
	node
	Keyword string `@("aggregate" | "Aggregat")`
}

type Aggregate struct {
	node
	KeywordAggregate *KeywordAggregate `@@`
	Name             *Name             `@@`
	Types            []*TypeDefinition `( "{" @@* "}" )?`
}

func (n *Aggregate) GetName() *Name {
	return n.Name
}

func (n *Aggregate) Children() []Node {
	res := sliceOf(n.KeywordAggregate, n.Name)
	for _, field := range n.Types {
		res = append(res, field)
	}

	return res
}

func (*Aggregate) namedType() {}
