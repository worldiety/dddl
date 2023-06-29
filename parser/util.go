package parser

import "reflect"

func sliceOf(nodes ...Node) []Node {
	var dst []Node
	for _, n := range nodes {
		if !isNil(n) {
			dst = append(dst, n)
		}
	}

	return dst
}

func isNil(a any) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
