package lsp

import (
	"fmt"
	"github.com/worldiety/dddl/lsp/protocol"
	p2 "github.com/worldiety/dddl/parser"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"strings"
)

// DYML language server.
type Server struct {
	// Map from Uri's to files.
	files map[protocol.DocumentURI]File
}

func NewServer() Server {
	fmt.Fprintln(os.Stderr, "hello world")
	return Server{
		files: make(map[protocol.DocumentURI]File),
	}
}

// Handle a client's request to initialize and respond with our capabilities.
func (s *Server) Initialize(params *protocol.InitializeParams) protocol.InitializeResult {
	log.Println(params)
	// For a perfect server we would need to check params.Capabilities to know
	// what information the client can handle.
	return protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.Full,
			SemanticTokensProvider: protocol.SemanticTokensOptions{
				Legend: protocol.SemanticTokensLegend{
					TokenTypes: TokenTypes,
				},
				Full: true,
			},
			CodeLensProvider: protocol.CodeLensOptions{
				ResolveProvider:         true,
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
			},
			HoverProvider: true,
		},
	}
}

// Initialized tells us, that the client is ready.
func (s *Server) Initialized() {
}

// Handle a hover event.
func (t *Server) Hover(params *protocol.HoverParams) protocol.Hover {
	return protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "markdown",
			Value: "## WTF?\n\n_IS_ this working?",
		},
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      1,
				Character: 1,
			},
			End: protocol.Position{
				Line:      1,
				Character: 3,
			},
		},
	} // Don't forget to enable hover capabilities when using this.
}

// A document was saved.
func (s *Server) DidSaveTextDocument(params *protocol.DidSaveTextDocumentParams) {
	s.sendDiagnostics()
}

// A document was opened.
func (s *Server) DidOpenTextDocument(params *protocol.DidOpenTextDocumentParams) {
	s.files[params.TextDocument.URI] = File{
		Uri:     params.TextDocument.URI,
		Content: params.TextDocument.Text,
	}
	s.sendDiagnostics()
}

// A document was close.
func (s *Server) DidCloseTextDocument(params *protocol.DidCloseTextDocumentParams) {
	delete(s.files, params.TextDocument.URI)
}

// A document was changed.
func (s *Server) DidChangeTextDocument(params *protocol.DidChangeTextDocumentParams) {
	// There is only a ever single full content change, as we requested.
	s.files[params.TextDocument.URI] = File{
		Uri:     params.TextDocument.URI,
		Content: params.ContentChanges[0].Text,
	}
	s.sendDiagnostics()
}

func (s *Server) FullSemanticTokens(params *protocol.SemanticTokensParams) protocol.SemanticTokens {
	// Mark "let" as a keyword for testing purposes.
	//fmt.Println("fullsemantic tokens", params)

	/*	file := s.files[params.TextDocument.URI]

		var data []uint32

		lexer := token.NewLexer(string(file.Uri), strings.NewReader(file.Content))

		// When a comment occurs the lexer emits a comment token and a chardata token.
		// We want to change the type of the chardata to be shown as a comment.
		nextCharIsComment := false

		for {
			tok, err := lexer.Token()
			if err != nil {
				// TODO What should we do here?
				break
			}

			part := SerializeToken(tok, nextCharIsComment)
			nextCharIsComment = false

			switch tok.Type() {
			case token.TokenG1Comment, token.TokenG2Comment:
				nextCharIsComment = true
			}

			data = append(data, part...)
		}

		// Make token positions relative.
		// Tokens are always 5 ints, first entry is line, second is char.
		for i := len(data) - 5; i >= 5; i -= 5 {
			// Make line difference relativ to previous
			data[i] -= data[i-5]
			// If item is in the same line, make char difference relative to previous
			if data[i] == 0 {
				data[i+1] -= data[i-5+1]
			}
		}*/

	log.Printf("%+v\n", params)
	//file := s.files[params.TextDocument.URI]
	doc, err := p2.Parse(string(params.TextDocument.URI[7:]))
	if err != nil {
		log.Println("cannot parse", err)
		return protocol.SemanticTokens{}
	}

	var tokens VSCTokens
	err = p2.Walk(doc, func(n p2.Node) error {
		// 1:3 -> 1:5 => just start and end col
		// 1:3 -> 2:5 => start until EOL and end from SOL to end col
		// 1:3 -> 3:5 => like above, but with full lines between

		start := n.Position()
		end := n.EndPosition()
		if start == end {
			log.Printf("token %T has invalid start/end: %+v->%+v\n", n, start, end)
			return nil // the token has not a useful token info attached
		}

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

	tokens = tokens[:]
	log.Printf("%+v\n", tokens)
	log.Printf("%+v\n", tokens.Encode())

	return protocol.SemanticTokens{
		Data: tokens.Encode(),
	}
}

func (s *Server) EncodeXML(filename protocol.DocumentURI) string {
	var out strings.Builder
	in := strings.NewReader(s.files[filename].Content)

	/*enc := encoder.NewXMLEncoder(filepath.Base(string(filename)), in, &out)
	err := enc.Encode()
	if err != nil {
		return ""
	}*/
	_ = in

	return out.String()
}

// sendDiagnostics sends any parser errors.
func (s *Server) sendDiagnostics() {
	/*
		for _, file := range s.files {

			fileContent := file.Content
			fileContent = strings.ToLower(fileContent)
			fileName := filepath.Base(string(file.Uri))

			// Parse file for any errors. Ideally we would be able to catch multiple errors and then recover.
			// Currently only the first error will be reported.
			diagnostics := []protocol.Diagnostic{}
			parser := parser.NewParser(fileName, strings.NewReader(fileContent))
			_, err := parser.Parse()
			if err != nil {
				switch e := err.(type) {
				case *token.PosError:
					for _, detail := range e.Details {
						diagnostics = append(diagnostics, protocol.Diagnostic{
							Range: protocol.Range{
								Start: protocol.Position{
									// Subtract 1 since dyml has 1 based lines and columns, but LSP wants 0 based
									Line:      uint32(detail.Node.Begin().Line) - 1,
									Character: uint32(detail.Node.Begin().Col) - 1,
								},
								End: protocol.Position{
									Line:      uint32(detail.Node.End().Line) - 1,
									Character: uint32(detail.Node.End().Col) - 1,
								},
							},
							Severity: protocol.SeverityError,
							Message:  e.Error(),
						})
					}
				default:
					diagnostics = append(diagnostics, protocol.Diagnostic{
						Severity: protocol.SeverityError,
						Message:  e.Error(),
					})
				}
			}

			_ = SendNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
				URI:         protocol.DocumentURI(file.Uri),
				Diagnostics: diagnostics,
			})
		}

	*/

	diag := protocol.Diagnostic{
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 0,
			},
			End: protocol.Position{
				Line:      0,
				Character: 3,
			},
		},
		Severity:           0,
		Code:               nil,
		CodeDescription:    nil,
		Source:             "",
		Message:            "vscode fuck off",
		Tags:               nil,
		RelatedInformation: nil,
		Data:               nil,
	}

	SendNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
		URI:         protocol.DocumentURI("file:///Users/tschinke/Downloads/nla/test.dyml.wdyspec"),
		Diagnostics: []protocol.Diagnostic{diag},
	})
}
