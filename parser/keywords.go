package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

type KeywordType struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("type" | "Typ")`
}

func (n *KeywordType) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordView struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("view" | "Ansicht")`
}

func (n *KeywordView) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordOutput struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("output" | "Ausgabe")`
}

func (n *KeywordOutput) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordInput struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("input" | "Eingabe")`
}

func (n *KeywordInput) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordWhile struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("while" | "solange")`
}

func (n *KeywordWhile) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordEvent struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("event" | "Ereignis")`
}

func (n *KeywordEvent) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordEventSent struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("sent" | "Zwischenereignis")`
}

func (n *KeywordEventSent) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordActivity struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("step" | "Aufgabe")`
}

func (n *KeywordActivity) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordWorkflow struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("workflow" | "Arbeitsablauf")`
}

func (n *KeywordWorkflow) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordDecision struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("decision" | "Entscheidung")`
}

func (n *KeywordDecision) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordIf struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("if" | "wenn")`
}

func (n *KeywordIf) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordElse struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("else" | "sonst")`
}

func (n *KeywordElse) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordActor struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("actor" | "Akteur")`
}

func (n *KeywordActor) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordReturn struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("return" | "Endereignis")`
}

func (n *KeywordReturn) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordReturnError struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("error" | "Fehler")`
}

func (n *KeywordReturnError) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordContext struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("context" | "Kontext")`
}

func (n *KeywordContext) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordTodo struct {
	node
	Tokens  []lexer.Token
	Keyword string `@"TODO"`
}

func (n *KeywordTodo) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}
