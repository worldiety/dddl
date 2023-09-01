package lsp

import (
	"fmt"
	"github.com/worldiety/dddl/lsp/protocol"
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/slices"
	"log"
)

// TokenTypes is a list of our supported types from the LSP spec.
// This array is sent once to the editor and after that only integers are used to refer
// to this array. See also
// https://code.visualstudio.com/api/language-extensions/semantic-highlight-guide
var TokenTypes = []string{"type", "string", "comment", "keyword", "struct", "function", "namespace"}

// These are indices into the TokenTypes array.
const (
	TokenType = iota
	TokenString
	TokenComment
	TokenKeyword
	TokenStruct
	TokenFunc
	TokenNamespace
)

// File is a file that is located at an Uri and has Content.
type File struct {
	Uri     protocol.DocumentURI
	Content string
}

type VSCToken struct {
	Node                                              parser.Node
	Line, StartChar, Length, TokenType, TokenModifier int
}

func (t VSCToken) String() string {
	return fmt.Sprintf("{%T: line: %d, col: %d, len: %d}\n", t.Node, t.Line, t.StartChar, t.Length)
}

type VSCTokens []VSCToken

func (v VSCTokens) FindBy(position protocol.Position) *VSCToken {
	for _, vscToken := range v {
		if vscToken.Line == int(position.Line) {
			if vscToken.StartChar <= int(position.Character) && (vscToken.StartChar+vscToken.Length) >= int(position.Character) {
				return &vscToken
			} else {
				log.Println("line ok, but char not matching", position.Character, vscToken.StartChar, vscToken.Length)
			}
		}
	}

	return nil
}

// Encode emits the delta encoded semantic tokens, see also
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_semanticTokens
// Note, that every 5 elements form a tuple (line, col, length, type, modifiers).
func (v VSCTokens) Encode() []uint32 {
	var res []uint32
	lastToken := VSCToken{Line: -1}
	for _, current := range v {
		var deltaLine, deltaStartChar uint32
		if lastToken.Line == -1 {
			deltaLine = uint32(current.Line)
			deltaStartChar = uint32(current.StartChar)
		} else {
			deltaLine = uint32(current.Line - lastToken.Line)
			if deltaLine == 0 {
				deltaStartChar = uint32(current.StartChar - lastToken.StartChar)
			} else {
				deltaStartChar = uint32(current.StartChar)
			}

		}

		lastToken = current
		res = append(res, deltaLine, deltaStartChar, uint32(current.Length), uint32(current.TokenType), uint32(current.TokenModifier))
	}

	return res
}

func getTokenType(node parser.Node) int {
	switch node.(type) {

	case *parser.KeywordTodo,
		*parser.KeywordContext,
		*parser.KeywordActor,
		*parser.KeywordEvent,
		*parser.KeywordEventSent,
		*parser.KeywordActivity,
		*parser.KeywordIf,
		*parser.KeywordElse,
		*parser.KeywordDecision,
		*parser.KeywordReturn,
		*parser.KeywordReturnError,
		*parser.KeywordWhile,
		*parser.KeywordView,
		*parser.KeywordInput,
		*parser.KeywordOutput,
		*parser.KeywordWorkflow:
		return TokenKeyword

	case *parser.Literal:
		return TokenString

	default:
		return TokenComment
	}
}

func isSemanticToken(n parser.Node) bool {
	switch n.(type) {
	case *parser.KeywordTodo,
		*parser.KeywordActor,
		*parser.KeywordEvent,
		*parser.KeywordEventSent,
		*parser.KeywordContext,
		*parser.KeywordActivity,
		*parser.KeywordIf,
		*parser.KeywordElse,
		*parser.KeywordDecision,
		*parser.KeywordReturn,
		*parser.KeywordReturnError,
		*parser.KeywordWhile,
		*parser.KeywordView,
		*parser.KeywordInput,
		*parser.KeywordOutput,
		*parser.KeywordWorkflow,
		*parser.KeywordChoice,
		*parser.KeywordFn,
		*parser.KeywordType,
		*parser.KeywordAlias,
		*parser.KeywordStruct:

		return true
	case *parser.Literal:
		return true
	case *parser.TypeDeclaration:
		return true
	case *parser.Choice, *parser.Type, *parser.Alias, *parser.Struct, *parser.Function:
		return true
	case parser.FnStmt:
		return true
	}

	return false
}

func IntoTokens(doc *parser.Doc) VSCTokens {
	var tokens VSCTokens
	err := parser.Walk(doc, func(n parser.Node) error {
		if !isSemanticToken(n) {
			return nil
		}
		// 1:3 -> 1:5 => just start and end col
		// 1:3 -> 2:5 => start until EOL and end from SOL to end col
		// 1:3 -> 3:5 => like above, but with full lines between

		start := n.Position()
		end := n.EndPosition()
		if start == end {
			log.Printf("token %T has invalid start/end: %+v->%+v\n", n, start, end)
			return nil // the token has not a useful token info attached
		}

		log.Printf("%T->%d:%d bis %d:%d\n", n, n.Position().Line, n.Position().Column, n.EndPosition().Line, n.EndPosition().Column)

		if start.Line == end.Line {
			tokens = append(tokens, VSCToken{
				Node:          n,
				Line:          start.Line - 1,
				StartChar:     start.Column - 1,
				Length:        end.Column - start.Column,
				TokenType:     getTokenType(n),
				TokenModifier: 0,
			})

			return nil
		} else {
			/*
				log.Printf("ignored: multiline token %T: %+v->%+v\n", n, start, end)

				tokens = append(tokens, VSCToken{
					Node:          n,
					Line:          start.Line - 1,
					StartChar:     start.Column - 1,
					Length:        1000, // don't know how long a line is
					TokenType:     getTokenType(n),
					TokenModifier: 0,
				})

				// everything in-between
				for i := 0; i < end.Line-start.Line; i++ {
					tokens = append(tokens, VSCToken{
						Node:          n,
						Line:          start.Line + i,
						StartChar:     0,    // don't know start-of-line
						Length:        1000, // don't know end-of-line
						TokenType:     getTokenType(n),
						TokenModifier: 0,
					})
				}

				tokens = append(tokens, VSCToken{
					Node:          n,
					Line:          end.Line - 1,
					StartChar:     0,
					Length:        end.Column, // don't know how long a line is
					TokenType:     getTokenType(n),
					TokenModifier: 0,
				})
			*/

		}

		/*


		 */

		// we don't have a length
		return nil
	})

	slices.SortFunc(tokens, func(a, b VSCToken) bool {
		if a.Line != b.Line {
			return a.Line < b.Line
		}

		return a.StartChar < b.StartChar
	})

	if err != nil {
		log.Println(err)
	}

	log.Println(tokens)

	return tokens
}
