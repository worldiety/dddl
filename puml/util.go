package puml

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
)

func TypeDeclToStr(decl *parser.TypeDeclaration) string {
	qname := resolver.NewQualifiedNameFromLocalName(decl.Name)
	tmp := qname.Name()
	if len(decl.Params) > 0 {
		tmp += "["
		for i, param := range decl.Params {
			tmp += TypeDeclToStr(param)
			if i < len(decl.Params)-1 {
				tmp += ", "
			}
		}
		tmp += "]"
	}

	return tmp
}

func typeDeclToLinkStr(r *resolver.Resolver, decl *parser.TypeDeclaration) string {
	qname := resolver.NewQualifiedNameFromLocalName(decl.Name)
	tmp := qname.Name()
	if len(decl.Params) > 0 {
		tmp += "[ " // significant whitespace, otherwise plantuml will omit the link
		for i, param := range decl.Params {
			tmp += typeDeclToLinkStr(r, param)
			if i < len(decl.Params)-1 {
				tmp += ", "
			}
		}
		tmp += " ]" // significant whitespace, otherwise plantuml will omit the link
	} else {
		defs := r.Resolve(qname)
		if len(defs) > 0 {
			// TODO this does not work properly in vsc, see also https://github.com/doxygen/doxygen/issues/7421
			tmp = "[[#" + qname.String() + " " + qname.Name() + "]]"
		}

	}

	return tmp
}

func record2Str(data *parser.Struct) string {
	if len(data.Fields) == 0 {
		return "Es wurden noch keine\nFelder definiert."
	}

	tmp := data.Name.Value + " = \n"
	for i, declaration := range data.Fields {
		tmp += TypeDeclToStr(declaration)
		if i < len(data.Fields)-1 {
			tmp += "\nund "
		}
	}

	return tmp
}

func choice2Str(data *parser.Choice) string {
	if len(data.Choices) == 0 {
		return "Es wurden noch keine\nWahlmöglichkeiten definiert."
	}

	tmp := data.Name.Value + " = \n"
	for i, declaration := range data.Choices {
		tmp += TypeDeclToStr(declaration)
		if i < len(data.Choices)-1 {
			tmp += "\noder "
		}
	}

	return tmp
}
