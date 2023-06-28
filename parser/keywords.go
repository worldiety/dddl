package parser

import "github.com/alecthomas/participle/v2/lexer"

type KeywordEvent struct {
	node
	Keyword string `("Event" | "Ereignis")`
}

func (n *KeywordEvent) Children() []Node {
	return nil
}

func (n *KeywordEvent) EndPosition() lexer.Position {
	return offsetPosText(n.Position(), n.Keyword)
}

type KeywordActivity struct {
	node
	Keyword string `("Step" | "Schritt")`
}

func (n *KeywordActivity) Children() []Node {
	return nil
}

func (n *KeywordActivity) EndPosition() lexer.Position {
	return offsetPosText(n.Position(), n.Keyword)
}

type KeywordActor struct {
	node
	Keyword string `("Actor" | "Akteur")`
}

func (n *KeywordActor) Children() []Node {
	return nil
}

func (n *KeywordActor) EndPosition() lexer.Position {
	return offsetPosText(n.Position(), n.Keyword)
}
