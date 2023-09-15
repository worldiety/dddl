package resolver

import (
	"github.com/worldiety/dddl/parser"
	"strings"
)

type FullQualifiedName string

func (n FullQualifiedName) String() string {
	return string(n)
}

func (n FullQualifiedName) Parent() FullQualifiedName {
	idx := strings.LastIndex(string(n), ".")
	if idx < 0 {
		return ""
	}

	return n[:idx]
}

func (n FullQualifiedName) Name() string {
	idx := strings.LastIndex(string(n), ".")
	if idx < 0 {
		return string(n)
	}

	return string(n[idx+1:])
}

func (n FullQualifiedName) With(name string) FullQualifiedName {
	if n == "" {
		return FullQualifiedName(name)
	}

	return n + "." + FullQualifiedName(name)
}

func NewQualifiedNameFromNamedType(n parser.NamedType) FullQualifiedName {
	name := n.GetName().Value
	parent := n.Parent()
	for parent != nil {
		if n, ok := parent.(parser.NamedType); ok {
			segment := secureSegment(n.GetName().Value)
			if segment == "" {
				return ""
			}
			name = segment + "." + name
		}
		parent = parent.Parent()
	}

	return FullQualifiedName(name)
}

func NewQualifiedNameFromLocalName(n *parser.QualifiedName) FullQualifiedName {
	name := ""
	for i, v := range n.Names {
		segment := secureSegment(v.Value)
		if segment == "" {
			return ""
		}
		name += segment
		if i < len(n.Names)-1 {
			name += "."
		}
	}

	if name == "" {
		return ""
	}

	if len(n.Names) > 1 {
		// by definition this must be already absolute
		return FullQualifiedName(name)
	}

	fqn := FullQualifiedName(name)
	child := n.Parent()

	for child.Parent() != nil {
		switch child.Parent().(type) {
		case *parser.Context, *parser.Aggregate:
			namedType := child.Parent().(parser.NamedType)
			fqn = FullQualifiedName(namedType.GetName().Value) + "." + fqn
		}

		child = child.Parent()
	}

	return fqn
}

func secureSegment(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, ".", "-"))
}

// NamedTypeFrom either returns the nearest parent which is a named type or nil.
func NamedTypeFrom(n parser.Node) parser.NamedType {
	parent := n
	for parent != nil {
		if nt, ok := parent.(parser.NamedType); ok {
			return nt
		}
		parent = parent.Parent()
	}

	return nil
}
