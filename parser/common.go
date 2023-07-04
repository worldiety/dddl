package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type Declaration interface {
	DeclaredName() *Ident
}

type Qualifier struct {
	Context *Context // can be nil for "shared kernel"
	Name    *Ident
}

func (q Qualifier) String() string {
	if q.Context == nil {
		return q.Name.String()
	}

	return q.Context.Name.String() + "." + q.Name.String()
}

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

func (n *Ident) String() string {
	if n == nil {
		return "nil"
	}

	return n.Value
}

func (n *Ident) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

func (n *Ident) IsUniverse() bool {
	// Boolean is not defined, we encourage to use a choicetype
	switch n.Value {
	//list, set, map, string, int, float
	case "Liste", "Menge", "Zuordnung", "Text", "Ganzzahl", "Zahl", "Gleitkommazahl":
		return true
	case "List", "Set", "Map", "String", "Number", "Integer", "Float":
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
	KeywordTodo *KeywordTodo `@@ ":"`
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
