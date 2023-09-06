package parser

import (
	"fmt"
	"strings"
)

type KeyValue struct {
	node
	Key   *Name `@@ ( "="`
	Value *Name `@@ )?`
}

func (n *KeyValue) String() string {
	if n == nil {
		return "nil"
	}

	key := ""
	if n.Key != nil {
		key = n.Key.Value
	}

	val := ""
	if n.Value != nil {
		val = n.Value.Value
	}

	return key + "=" + val
}

func (n *KeyValue) Children() []Node {
	return sliceOf(n.Key, n.Value)
}

type Annotation struct {
	node
	Name      *Name       `@@`
	KeyValues []*KeyValue `( "("  @@ ("," @@)*  ")" )?`
}

func (n *Annotation) ExpectEmpty() error {
	if len(n.KeyValues) != 0 {
		var tmp []string
		for _, value := range n.KeyValues {
			tmp = append(tmp, value.String())
		}
		return fmt.Errorf("expected empty annotation map, but got (%s)", strings.Join(tmp, ","))
	}

	return nil
}

func (n *Annotation) ExpectKeysOf(possibleKeySet ...string) error {
	tmp := map[string]bool{}
	for _, kv := range n.KeyValues {
		found := false
		for _, s := range possibleKeySet {
			if s == kv.Key.Value {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unexpected key '%s', must be of (%s)", kv.Key.Value, strings.Join(possibleKeySet, "|"))
		}

		if tmp[kv.Key.Value] {
			return fmt.Errorf("key has already been defined: '%s'", kv.Key.Value)
		}

		tmp[kv.Key.Value] = true
	}

	return nil
}

func (n *Annotation) FirstValue(possibleKeySet ...string) (value string, found bool) {
	for _, kv := range n.KeyValues {
		for _, s := range possibleKeySet {
			if s == kv.Key.Value {
				if kv.Value == nil {
					return "", true
				}

				return kv.Value.Value, true
			}
		}

	}

	return "", false
}

func (n *Annotation) Children() []Node {
	res := sliceOf(n.Name)
	for _, value := range n.KeyValues {
		res = append(res, value)
	}

	return res
}

type TypeDefinition struct {
	node
	// Description may be nil
	Description *Literal      `@@?`
	Annotations []*Annotation `("@" @@)*`
	Type        NamedType     `@@`
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
