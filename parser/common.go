package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
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

// Ident refers to the rules of an Identifier used by the Lexer.
type Ident struct {
	node
	Tokens []lexer.Token
	Value  string `@Name`
}

func NewIdentWithParent(parent Node, value string) *Ident {
	return &Ident{
		node: node{
			Pos:    parent.Position(),
			EndPos: parent.EndPosition(),
			parent: parent,
		},
		Value: value,
	}
}

func (n *Ident) String() string {
	if n == nil {
		return "nil"
	}

	return n.Value
}

func (n *Ident) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

func (n *Ident) IsList() bool {
	return n.Value == UList || n.Value == UListDE
}

func (n *Ident) IsSet() bool {
	return n.Value == USet || n.Value == USetDE
}

func (n *Ident) IsMap() bool {
	return n.Value == UMap || n.Value == UMapDE
}

func (n *Ident) IsNumber() bool {
	return n.Value == UNumber || n.Value == UNumberDE
}

func (n *Ident) IsFloat() bool {
	return n.Value == UFloat || n.Value == UFloatDE
}

func (n *Ident) IsInt() bool {
	return n.Value == UInt || n.Value == UIntDE
}

func (n *Ident) IsAny() bool {
	return n.Value == UAny
}

func (n *Ident) IsFunc() bool {
	return n.Value == UFunc
}

func (n *Ident) IsString() bool {
	return n.Value == UString || n.Value == UStringDE
}

func (n *Ident) NormalizeUniverse() string {
	if n.IsList() {
		return UList
	}

	if n.IsSet() {
		return USet
	}

	if n.IsMap() {
		return UMap
	}

	if n.IsNumber() {
		return UNumber
	}

	if n.IsFloat() {
		return UFloat
	}

	if n.IsInt() {
		return UInt
	}

	if n.IsAny() {
		return UAny
	}

	if n.IsFunc() {
		return UFunc
	}

	if n.IsString() {
		return UString
	}

	return n.Value
}

const (
	UList     = "List"
	UListDE   = "Liste"
	USet      = "Set"
	USetDE    = "Menge"
	UMap      = "Map"
	UMapDE    = "Zuordnung"
	UString   = "String"
	UStringDE = "Text"
	UNumber   = "Number"
	UNumberDE = "Zahl"
	UInt      = "Integer"
	UIntDE    = "Ganzzahl"
	UFloat    = "Float"
	UFloatDE  = "Gleitkommazahl"
	UAny      = "any"
	UFunc     = "func"
	UError    = "error"
	UContext  = "context"
)

func (n *Ident) IsUniverse() bool {
	// Boolean is not defined, we encourage to use a choicetype
	switch n.Value {
	//list, set, map, string, int, float
	case UListDE, USetDE, UMapDE, UStringDE, UNumberDE, UIntDE, UFloatDE:
		return true
	case UList, USet, UMap, UString, UNumber, UInt, UFloat:
		return true
	case UAny, UFunc, UError, UContext:
		return true
	default:
		return false
	}
}

type Definition struct {
	node
	Tokens []lexer.Token
	Text   string `@Text`
}

func (n *Definition) Value() string {
	if n == nil {
		return ""
	}

	return n.Text
}

func (n *Definition) EndPosition() lexer.Position {
	pos := n.relocateEndPos(n.Tokens)
	pos.Column += 2 // fix leading and appended "
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

type IdentOrLiteral struct {
	node
	Name    *Ident   `(@@`
	Literal *Literal `|@@)`
}

func (n *IdentOrLiteral) EndPosition() lexer.Position {
	if n.Name != nil {
		return n.Name.EndPosition()
	}

	return n.Literal.EndPosition()
}

// Value returns either the Names' value or the Literals' value.
func (n *IdentOrLiteral) Value() string {
	if n.Name != nil {
		return n.Name.Value
	}

	return n.Literal.Value
}

func (n *IdentOrLiteral) Children() []Node {
	return sliceOf(n.Name, n.Literal)
}

type ToDo struct {
	node
	KeywordTodo *KeywordTodo `@@`
	Text        *ToDoText    `@@`
}

func (n *ToDo) Value() string {
	if n == nil || n.Text == nil {
		return ""
	}

	return n.Text.Text
}

func (n *ToDo) Children() []Node {
	return sliceOf(n.KeywordTodo, n.Text)
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

type ToDoText struct {
	node
	Tokens []lexer.Token
	Text   string `@Text`
}

func (n *ToDoText) Children() []Node {
	return nil
}

func (n *ToDoText) EndPosition() lexer.Position {
	pos := n.relocateEndPos(n.Tokens)
	pos.Column += 2 // fix leading and appended "
	return pos
}
