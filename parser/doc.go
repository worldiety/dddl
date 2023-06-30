package parser

type ContextOrTypeDefinition struct {
	node
	Context        *Context        `@@`
	TypeDefinition *TypeDefinition `|@@`
}

func (n *ContextOrTypeDefinition) Children() []Node {
	return sliceOf(n.Context, n.TypeDefinition)
}

type Doc struct {
	node
	Definitions []*ContextOrTypeDefinition `@@*`
}

func (n *Doc) ContextByName(name string) *Context {
	for _, definition := range n.Definitions {
		if context := definition.Context; context != nil {
			if context.Name.Value == name {
				return context
			}
		}
	}

	return nil
}

func (n *Doc) Contexts() []*Context {
	var res []*Context
	for _, definition := range n.Definitions {
		if context := definition.Context; context != nil {
			res = append(res, context)
		}
	}

	return res
}

func (n *Doc) DataByName(name string) *Data {
	for _, definition := range n.Definitions {
		if context := definition.Context; context != nil {
			for _, element := range context.Elements {
				if element.DataType != nil {
					if element.DataType.Name.Value == name {
						return element.DataType
					}
				}
			}
		}
	}

	return nil
}

func (n *Doc) WorkflowByName(name string) *Workflow {
	for _, definition := range n.Definitions {
		if context := definition.Context; context != nil {
			for _, element := range context.Elements {
				if element.Workflow != nil {
					if element.Workflow.Name.Value == name {
						return element.Workflow
					}
				}
			}
		}
	}

	return nil
}

func (n *Doc) Children() []Node {
	var res []Node
	for _, context := range n.Definitions {
		res = append(res, context)
	}

	return res
}
