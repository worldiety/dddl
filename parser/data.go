package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

// Ident refers to the rules of an Identifier used by the Lexer.
type Ident struct {
	node

	Name string `@Name`
}

func (n *Ident) EndPosition() lexer.Position {
	return offsetPosText(n.Position(), n.Name)
}

func (n *Ident) Children() []Node {
	return nil
}

func (n *Ident) IsUniverse() bool {
	// Boolean is not defined, we encourage to use a choicetype
	switch n.Name {
	//list, set, map, string, int, float
	case "Liste", "Menge", "Zuordnung", "Text", "Zahl", "Gleitkommazahl":
		return true
	default:
		return false
	}
}

type KeywordTodo struct {
	node
	Keyword string `"TODO"`
}

func (d *KeywordTodo) Children() []Node {
	return nil
}

func (n *KeywordTodo) EndPosition() lexer.Position {
	return offsetPosText(n.Position(), "TODO")
}

type ToDo struct {
	node
	KeywordTodo *KeywordTodo `@@ ":"`
	Text        *ToDoText    `@@`
}

func (n *ToDo) Children() []Node {
	return []Node{n.KeywordTodo, n.Text}
}

func offsetPosText(pos lexer.Position, text string) lexer.Position {
	if len(text) == 0 {
		return pos
	}

	lines := strings.Split(text, "\n")
	var lastLineLen int
	if len(lines) == 1 {
		lastLineLen = pos.Column + len(lines[len(lines)-1])
	} else {
		lastLineLen = len(lines[len(lines)-1])
	}

	pos.Line += len(lines) - 1
	pos.Column = lastLineLen

	return pos
}

type Definition struct {
	node
	Text string `@Text`
}

func (n *Definition) Position() lexer.Position {
	pos := n.node.Position()
	pos.Column++ // fix "
	return pos
}

func (n *Definition) EndPosition() lexer.Position {
	pos := offsetPosText(n.Position(), n.Text)
	return pos
}

type Text struct {
	node
	Value string `@Text`
}

func (n *Text) Children() []Node {
	return nil
}

type ToDoText struct {
	node
	Text string `@Text`
}

func (n *ToDoText) Children() []Node {
	return nil
}

func (n *ToDoText) Position() lexer.Position {
	pos := n.node.Position()
	pos.Column++ // fix "
	return pos
}

func (n *ToDoText) EndPosition() lexer.Position {
	pos := offsetPosText(n.Position(), n.Text)
	return pos
}

func (n *Definition) Empty() bool {
	if n == nil {
		return true
	}

	tmp := strings.TrimSpace(n.Text)
	if tmp == "" {
		return true
	}

	if tmp == "???" {
		return true
	}

	return false
}

func (n *Definition) NeedsRevise() bool {
	if n.Empty() {
		return false
	}

	return strings.Contains(n.Text, "???")
}

func (n *Definition) Children() []Node {
	return nil
}

type KeywordData struct {
	node
	Keyword string `"Daten"`
}

func (d *KeywordData) Children() []Node {
	return nil
}

func (n *KeywordData) EndPosition() lexer.Position {
	return offsetPosText(n.Position(), "Daten")
}

type KeywordContext struct {
	node
	Keyword string `"Kontext"`
}

func (d *KeywordContext) Children() []Node {
	return nil
}

func (n *KeywordContext) EndPosition() lexer.Position {
	return offsetPosText(n.Position(), "Kontext")
}

// A Data is either a choice or a compound data type.
// Combining both is probably hard to understand.
// Without massive lookahead, we cannot distinguish that, so we will
// check that using a linter later.
type Data struct {
	node
	KeywordData *KeywordData       `@@`
	Name        *Ident             ` @@ ( "="`
	ToDo        *ToDo              `@@? `
	First       *TypeDeclaration   ` (@@ `
	Fields      []*TypeDeclaration `("und" @@)*`
	Choices     []*TypeDeclaration `("oder" @@)*)?  `

	Definition *Definition `@@?)?`
}

func (d *Data) Empty() bool {
	return d.First == nil && len(d.Fields) == 0 && len(d.Choices) == 0
}

func (d *Data) Children() []Node {
	var res []Node
	res = append(res, d.KeywordData, d.Name)

	for _, declaration := range d.ChoiceTypes() {
		res = append(res, declaration)
	}

	for _, declaration := range d.FieldTypes() {
		res = append(res, declaration)
	}

	if d.Definition != nil {
		res = append(res, d.Definition)
	}

	if d.ToDo != nil {
		res = append(res, d.ToDo)
	}

	return res
}

// ChoiceTypes is nil, if any field is defined.
// If neither choices nor fields are defined, this also returns nil.
// This avoids visiting nodes twice.
func (d *Data) ChoiceTypes() []*TypeDeclaration {
	if len(d.Fields) > 0 || len(d.Choices) == 0 {
		return nil
	}

	var choices []*TypeDeclaration
	choices = append(choices, d.First)
	choices = append(choices, d.Choices...)
	return choices
}

// FieldTypes is nil if any choice type is defined, otherwise contains
// at least the First declaration.
func (d *Data) FieldTypes() []*TypeDeclaration {
	if len(d.Choices) > 0 {
		return nil
	}

	var fields []*TypeDeclaration
	if d.First != nil {
		fields = append(fields, d.First)
	}
	fields = append(fields, d.Fields...)
	return fields
}
