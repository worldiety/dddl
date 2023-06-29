package parser

// Literal refers to the rules of quoted Text by the Lexer.
type Literal struct {
	node
	Value string `@Text`
}

// Ident refers to the rules of an Identifier used by the Lexer.
type Ident struct {
	node
	Value string `@Name`
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
	Text string `@Text`
}

type IdentOrLiteral struct {
	node
	Name    *Ident   `(@@`
	Literal *Literal `|@@)`
}

// Value returns either the Names' value or the Literals' value.
func (n *IdentOrLiteral) Value() string {
	if n.Name != nil {
		return n.Name.Value
	}

	return n.Literal.Value
}

func (n *IdentOrLiteral) Children() []Node {
	if n.Name != nil {
		return []Node{n.Name}
	}
	return []Node{n.Literal}
}
