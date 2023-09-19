package parser

type TypeDefinition struct {
	node
	// Description may be nil
	Description       *Literal      `@@?`
	Annotations       []*Annotation `("@" @@)*`
	Type              NamedType     `@@`
	parsedAnnotations []TypedAnnotation
}

func TypeDefinitionFrom(n Node) *TypeDefinition {
	for n != nil {
		if td, ok := n.(*TypeDefinition); ok {
			return td
		}

		n = n.Parent()
	}

	return nil
}

// TypedAnnotations parses the Annotations once and caches them internally.
func (n *TypeDefinition) TypedAnnotations() []TypedAnnotation {
	if n.parsedAnnotations == nil {
		n.parsedAnnotations = parseAnnotations(n)
	}

	return n.parsedAnnotations
}

func (n *TypeDefinition) Children() []Node {
	res := sliceOf(n.Description)
	for _, tag := range n.Annotations {
		res = append(res, tag)
	}

	res = append(res, n.Type)
	return res
}

type NamedType interface {
	Node
	namedType()
	GetName() *Name
	GetKeyword() string
}

func FindAnnotation[T TypedAnnotation](n Node) T {
	if tDef := TypeDefinitionFrom(n); tDef != nil {
		for _, annotation := range tDef.TypedAnnotations() {
			if a, ok := annotation.(T); ok {
				return a
			}
		}
	}

	var zero T
	return zero
}
