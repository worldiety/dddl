package linter

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
	"log"
)

type UndeclaredTypeDeclInNamedType struct {
	hint
	Parent   *parser.TypeDefinition
	TypeDecl *parser.TypeDeclaration
}

type FirstUndeclaredTypeDeclInNamedType struct {
	hint
	Parent   *parser.TypeDefinition
	TypeDecl *parser.TypeDeclaration
}

// CheckUndefined searches for all Identifiers and checks
// if there are workflows or data types for them.
func CheckUndefined(r *resolver.Resolver) []Hint {
	var res []Hint

	dedupTableNames := map[string]struct{}{}
	parser.MustWalk(r.Workspace(), func(n parser.Node) error {
		if decl, ok := n.(*parser.TypeDeclaration); ok {
			qname := resolver.NewQualifiedNameFromLocalName(decl.Name)
			defs := r.Resolve(qname)
			if len(defs) == 0 {
				namedType := parser.TypeDefinitionFrom(decl)
				if namedType == nil {
					log.Println("failed to get nearest type definition from declaration: ", decl.Name.String())
					return nil
				}

				parentQname := resolver.NewQualifiedNameFromNamedType(namedType.Type)
				if _, ok := dedupTableNames[parentQname.String()]; !ok {
					dedupTableNames[parentQname.String()] = struct{}{}
					res = append(res, &FirstUndeclaredTypeDeclInNamedType{
						Parent:   namedType,
						TypeDecl: decl,
					})
				}

				res = append(res, &UndeclaredTypeDeclInNamedType{
					Parent:   namedType,
					TypeDecl: decl,
				})
			}
		}
		return nil
	})

	return res
}
