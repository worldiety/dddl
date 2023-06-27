package lsp

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/lsp/protocol"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/web/editor"
	"golang.org/x/exp/slog"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type jobQueueEntry struct {
	lastId int
	queue  chan func()
}

// DYML language server.
type Server struct {
	// Map from Uri's to files.
	files             map[protocol.DocumentURI]File
	lastPreviewParams *PreviewHtmlParams
	jobQueues         map[string]jobQueueEntry
	jobQueuesLock     sync.Mutex
	rootPath          string
}

func NewServer() *Server {
	return &Server{
		files:     make(map[protocol.DocumentURI]File),
		jobQueues: map[string]jobQueueEntry{},
	}
}

func (s *Server) reloadFiles() {
	err := filepath.Walk(s.rootPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if strings.HasSuffix(path, ".ddd") {
			buf, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read '%s': %w", path, err)
			}

			uri := protocol.DocumentURI("file://" + path)
			s.files[uri] = File{
				Uri:     uri,
				Content: string(buf),
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("failed to reload files: %v\n", err)
	}
}

// Handle a client's request to initialize and respond with our capabilities.
func (s *Server) Initialize(params *protocol.InitializeParams) protocol.InitializeResult {
	log.Printf("%+v", params)
	s.rootPath = params.RootPath
	s.reloadFiles()

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
func (s *Server) Hover(params *protocol.HoverParams) protocol.Hover {
	file := s.files[params.TextDocument.URI]
	doc, err := parser.ParseText(string(file.Uri), file.Content)
	if err != nil {
		log.Println("cannot parse", err)
		return protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: "## Syntaxfehler\nPrüfe deinen Text und die Fehlermeldung: " + err.Error(),
			},
		}
	}

	tokens := IntoTokens(doc)
	token := tokens.FindBy(params.Position)
	if token == nil {
		log.Println("token not found")
		return protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: "",
			},
		}
	}

	return protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "markdown",
			Value: s.hoverText(token),
		},
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(token.Node.Position().Line - 1),
				Character: uint32(token.Node.Position().Column - 1),
			},
			End: protocol.Position{
				Line:      uint32(token.Node.EndPosition().Line - 1),
				Character: uint32(token.Node.EndPosition().Column - 1),
			},
		},
	}
}

func (s *Server) hoverText(token *VSCToken) string {
	switch token.Node.(type) {
	case *parser.Context:
		return "Schlüsselwort, dass einen _Bounded Context_ einführt"
	case *parser.ToDoText:
		return "Textliteral, dass eine noch zu erledigende Aufgabe beschreibt."
	case *parser.KeywordData:
		return "Schlüsselwort, dass eine bestimmte Art von Daten einführt."
	case *parser.Ident:
		return "Dies ist ein Bezeichner, der nur genau einmal eindeutig eingeführt darf."
	case *parser.Definition:
		return "Diese Definition gehört zu einem Bezeichner und beschreibt den Sachverhalt ausführlich."
	default:
		return fmt.Sprintf("%T", token.Node)
	}

}

// A document was saved.
func (s *Server) DidSaveTextDocument(params *protocol.DidSaveTextDocumentParams) {
	if params.Text != nil {
		// this happens by definition, but why and when exactly?
		s.files[params.TextDocument.URI] = File{
			Uri:     params.TextDocument.URI,
			Content: *params.Text,
		}
	}

	s.sendDiagnostics()
	s.sendPreviewHtml()
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
	s.sendPreviewHtml()
}

func (s *Server) FullSemanticTokens(params *protocol.SemanticTokensParams) protocol.SemanticTokens {
	file := s.files[params.TextDocument.URI]
	doc, err := parser.ParseText(string(file.Uri), file.Content)
	if err != nil {
		log.Println("cannot parse", err)
		return protocol.SemanticTokens{}
	}

	tokens := IntoTokens(doc)

	tokens = tokens[:]
	//log.Printf("%+v\n", tokens)
	//log.Printf("%+v\n", tokens.Encode())

	return protocol.SemanticTokens{
		Data: tokens.Encode(),
	}
}

func (s *Server) AsciiDoc(filename protocol.DocumentURI) string {
	var out strings.Builder
	doc, err := s.parseSuperDoc()
	if doc == nil {
		return err.Error()
	}

	out.WriteString("= Implement me\n\n")
	for _, context := range doc.Contexts {
		out.WriteString("== ")
		out.WriteString(context.Name.Name)
		out.WriteString("\n")
	}

	return out.String()
}

func (s *Server) sendPreviewHtml() {
	s.async("previewHtml", func() {
		if s.lastPreviewParams == nil {
			slog.Warn("cannot send preview html, never requested a preview")
			err := SendNotification("custom/newAsyncPreviewHtml", "missing WebView preview params from LSP client")
			if err != nil {
				log.Printf("cannot send unbound previewhtml: %v", err)
			}
			return
		}

		html := s.RenderPreviewHtml(*s.lastPreviewParams)
		err := SendNotification("custom/newAsyncPreviewHtml", html)
		if err != nil {
			log.Printf("cannot send previewhtml: %v", err)
		}
	})
}

// sendDiagnostics sends any parser errors.
func (s *Server) sendDiagnostics() {
	s.async("diagnostics", func() {
		for _, file := range s.files {
			var diagnostics []protocol.Diagnostic

			_, err := parser.ParseText(string(file.Uri), file.Content)
			if err != nil {
				switch err := err.(type) {
				case participle.Error:
					pos := err.Position()

					diagnostics = append(diagnostics, protocol.Diagnostic{
						Range: protocol.Range{
							Start: protocol.Position{
								// Subtract 1 since dyml has 1 based lines and columns, but LSP wants 0 based
								Line:      uint32(pos.Line) - 1,
								Character: uint32(pos.Column) - 1,
							},
							// we don't know the length, so just always pick the next 3 chars
							End: protocol.Position{
								Line:      uint32(pos.Line) - 1,
								Character: uint32(pos.Column+3) - 1,
							},
						},
						Severity: protocol.SeverityError,
						Message:  err.Error(),
					})

				default:
					diagnostics = append(diagnostics, protocol.Diagnostic{
						Severity: protocol.SeverityError,
						Message:  err.Error(),
					})
				}
			} else {
				// we must always send the diagnostics, otherwise error message will not disappear
				diagnostics = make([]protocol.Diagnostic, 0)
			}

			err = SendNotification("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
				URI:         protocol.DocumentURI(file.Uri),
				Diagnostics: diagnostics,
			})

			if err != nil {
				log.Printf("cannot send diagnostics: %v", err)
			}

		}
	})

}

type PreviewHtmlParams struct {
	Doc         protocol.DocumentURI
	TailwindUri protocol.DocumentURI
}

func (s *Server) parseSuperDoc() (*parser.Doc, error) {
	var superErr error
	superDoc := &parser.Doc{}
	for _, file := range s.files {
		doc, err := parser.ParseText(string(file.Uri), file.Content)
		if err != nil {
			if superErr == nil {
				superErr = err
			} else {
				superErr = fmt.Errorf("also occurred: %w and %w", superErr, err)
			}
		}

		for _, context := range doc.Contexts {
			ctx := superDoc.ContextByName(context.Name.Name)
			if ctx == nil {
				superDoc.Contexts = append(superDoc.Contexts, context)
			} else {
				for _, element := range context.Elements {
					ctx.Elements = append(ctx.Elements, element) // just append it
				}

				if ctx.Definition == nil {
					ctx.Definition = context.Definition // maybe non-nil, TODO double definition is a lint-error
				}
			}
		}
	}

	return superDoc, superErr
}

func (s *Server) RenderPreviewHtml(params PreviewHtmlParams) string {
	log.Println(params)
	s.lastPreviewParams = &params

	var model editor.EditorPreview
	model.VSCode.ScriptUris = append(model.VSCode.ScriptUris, string(s.lastPreviewParams.TailwindUri))

	doc, err := s.parseSuperDoc()
	if doc == nil {
		return err.Error()
	}

	if err != nil {
		model.Error = err.Error()
	}

	linter := editor.Linter(func(doc *parser.Doc) []linter.Hint {
		return linter.Lint(doc)
	})

	return editor.RenderViewHtml(linter, doc, model)
}
