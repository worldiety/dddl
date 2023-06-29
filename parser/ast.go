package parser

import "github.com/alecthomas/participle/v2/lexer"

type Node interface {
	Position() lexer.Position
	EndPosition() lexer.Position
	Children() []Node // this simplifies the code style but has a lot of extra allocations
}

type node struct {
	Pos    lexer.Position
	EndPos lexer.Position
}

func (n *node) Position() lexer.Position {
	return n.Pos
}

func (n *node) EndPosition() lexer.Position {
	return n.EndPos
}

func (n *node) Children() []Node {
	return nil
}

func Walk(n Node, visitor func(n Node) error) error {
	if n == nil {
		return nil
	}

	if err := visitor(n); err != nil {
		return err
	}

	for _, c := range n.Children() {
		if err := Walk(c, visitor); err != nil {
			return err
		}
	}

	return nil
}
