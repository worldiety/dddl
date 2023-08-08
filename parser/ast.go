package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type Node interface {
	Position() lexer.Position
	EndPosition() lexer.Position
	Children() []Node // this simplifies the code style but has a lot of extra allocations
	Parent() Node
	setParent(p Node)
}

type node struct {
	Pos    lexer.Position
	EndPos lexer.Position
	parent Node
}

func (n *node) Parent() Node {
	return n.parent
}
func (n *node) setParent(p Node) {
	n.parent = p
}

func (n *node) Position() lexer.Position {
	return n.Pos
}

func (n *node) EndPosition() lexer.Position {
	return n.EndPos
}

// relocateEndPos uses the given lexer token slice to recalculate the actual
// endposition which is often wrong. This is a bug in the participle parser, which
// appends all whitespace until the next token appears.
func (n *node) relocateEndPos(tokens []lexer.Token) lexer.Position {
	if len(tokens) == 0 {
		return n.EndPos
	}

	pos := n.Position()
	pos.Column += len(strings.TrimSpace(tokens[len(tokens)-1].Value))
	return pos
}

func (n *node) Children() []Node {
	return nil
}

func MustWalk(n Node, visitor func(n Node) error) {
	if err := Walk(n, visitor); err != nil {
		panic(fmt.Errorf("unexpected %w", err))
	}
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

func WorkspaceOf(root Node) *Workspace {
	for root != nil {
		if ws, ok := root.(*Workspace); ok {
			return ws
		}

		root = root.Parent()
	}

	return nil
}
