package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

// Literal refers to the rules of quoted Text by the Lexer.
type Literal struct {
	node
	Tokens []lexer.Token
	Value  string `@Text`
}

func (n *Literal) EndPosition() lexer.Position {
	pos := n.relocateEndPos(n.Tokens)
	pos.Column += 2 // fix leading and appended "
	return pos
}

type Name struct {
	node
	Tokens []lexer.Token
	Value  string `(@Name | @Text)`
}

func (n *Name) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type QualifiedName struct {
	node
	Names []Name `@@ ("." @@)*`
}

func (n *QualifiedName) String() string {
	if len(n.Names) == 1 {
		return n.Names[0].Value
	}

	tmp := ""
	for i, name := range n.Names {
		tmp += name.Value
		if i < len(n.Names)-1 {
			tmp += "."
		}
	}
	return tmp
}
