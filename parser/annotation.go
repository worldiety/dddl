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

type TypedAnnotation interface {
	typedAnnotation()
}
