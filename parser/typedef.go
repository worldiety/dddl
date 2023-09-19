package parser

import (
	"fmt"
	"strings"
)

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
		n.parsedAnnotations = make([]TypedAnnotation, 0, len(n.Annotations)) // optimize and allocate on 0 size

	}

	return n.parsedAnnotations
}

func (n *TypeDefinition) ExpectOnlyOf(names ...string) error {
	tmp := map[string]int{}
	for _, annotation := range n.Annotations {
		tmp[annotation.Name.Value] = tmp[annotation.Name.Value] + 1
	}

	for k, count := range tmp {
		if count > 1 {
			return fmt.Errorf("key '%s' is ambigous", k)
		}

		allowed := false
		for _, name := range names {
			if name == k {
				allowed = true
				break
			}
		}

		if !allowed {
			return fmt.Errorf("key '%s' must not be defined on type '%s'", k, n.Type.GetName().Value)
		}
	}

	return nil
}

func (n *TypeDefinition) ExpectOneOrNoneOf(names ...string) (*Annotation, error) {
	var res []*Annotation
	for _, annotation := range n.Annotations {
		for _, name := range names {
			if annotation.Name.Value == name {
				res = append(res, annotation)
			}
		}

	}

	if len(res) == 0 {
		return nil, nil
	}

	if len(res) > 1 {
		var got []string
		for _, r := range res {
			got = append(got, r.Name.Value)
		}

		return nil, fmt.Errorf("expected none or exact one of (%s) but got (%s)", strings.Join(names, "|"), strings.Join(got, "+"))
	}

	return res[0], nil
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
