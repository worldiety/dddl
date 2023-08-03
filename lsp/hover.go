package lsp

import (
	_ "embed"
	"fmt"
	"github.com/worldiety/dddl/lsp/protocol"
	"github.com/worldiety/dddl/parser"
	"log"
	"strings"
)

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

//go:embed tips/kw_context_decl.md
var tipKWContextDecl string

//go:embed tips/kw_context_workflowdef.md
var tipKWContextWFDef string

//go:embed tips/declaredident.md
var tipDeclarationIdent string

//go:embed tips/definedident.md
var tipDefinedIdent string

//go:embed tips/universeident.md
var tipUniverseIdent string

//go:embed tips/definition.md
var tipDefinition string

//go:embed tips/todotext.md
var tipTodoText string

//go:embed tips/kw_data.md
var tipKWData string

//go:embed tips/kw_actor.md
var tipKWActor string

//go:embed tips/kw_task.md
var tipKWTask string

//go:embed tips/kw_workflow.md
var tipKWWorkflow string

//go:embed tips/kw_view.md
var tipKWView string

//go:embed tips/kw_input.md
var tipKWInput string

//go:embed tips/kw_output.md
var tipKWOutput string

//go:embed tips/kw_return.md
var tipKWReturn string

//go:embed tips/wf_actor.md
var tipWFActor string

//go:embed tips/wf_task.md
var tipWFTask string

//go:embed tips/wf_view.md
var tipWFView string

//go:embed tips/wf_input.md
var tipWFInput string

//go:embed tips/wf_output.md
var tipWFOutput string

//go:embed tips/wf_if.md
var tipWFIf string

//go:embed tips/wf_while.md
var tipWFWhile string

//go:embed tips/wf_returnerr.md
var tipWFReturnErr string

//go:embed tips/wf_eventsent.md
var tipWFEventSent string

func (s *Server) hoverText(token *VSCToken) string {
	switch n := token.Node.(type) {
	case *parser.KeywordContext:
		if _, isDecl := n.Parent().(*parser.Context); isDecl {
			return fmt.Sprintf(tipKWContextDecl, n.Keyword)
		} else {
			return fmt.Sprintf(tipKWContextWFDef, n.Keyword)
		}
	case *parser.Ident:
		declaration := false
		switch n.Parent().(type) {
		case *parser.Function, *parser.Alias, *parser.Struct, *parser.Choice, *parser.Type, *parser.Context:
			declaration = true
		}
		if declaration {
			return fmt.Sprintf(tipDeclarationIdent, n.Value)
		}

		if n.IsUniverse() {
			return fmt.Sprintf(tipUniverseIdent, n.Value)
		}

		return fmt.Sprintf(tipDefinedIdent, n.Value)

	case *parser.Definition:
		s := shortStringLit(n.Text)
		return fmt.Sprintf(tipDefinition, s)

	case *parser.ToDoText:
		s := shortStringLit(n.Text)
		return fmt.Sprintf(tipTodoText, s)
	case *parser.KeywordData:
		return fmt.Sprintf(tipKWData, n.Keyword)
	case *parser.KeywordWorkflow:
		return fmt.Sprintf(tipKWWorkflow, n.Keyword)
	case *parser.KeywordActor:
		return fmt.Sprintf(tipKWActor, n.Keyword)
	case *parser.KeywordActivity:
		return fmt.Sprintf(tipKWTask, n.Keyword)
	case *parser.KeywordView:
		return fmt.Sprintf(tipKWView, n.Keyword)
	case *parser.KeywordInput:
		return fmt.Sprintf(tipKWInput, n.Keyword)
	case *parser.KeywordOutput:
		return fmt.Sprintf(tipKWOutput, n.Keyword)
	case *parser.KeywordReturn:
		return fmt.Sprintf(tipKWReturn, n.Keyword)
	case *parser.Literal:
		if identOrLit, ok := n.Parent().(*parser.IdentOrLiteral); ok {
			switch identOrLit.Parent().(type) {
			/*case *parser.ActorStmt:
			return fmt.Sprintf(tipWFActor, n.Value)*/
			case *parser.Function:
				return fmt.Sprintf(tipWFTask, n.Value)
			/*case *parser.ViewStmt:
				return fmt.Sprintf(tipWFView, n.Value)
			case *parser.InputStmt:
				return fmt.Sprintf(tipWFInput, n.Value)
			case *parser.OutputStmt:
				return fmt.Sprintf(tipWFOutput, n.Value)
			case *parser.ReturnErrorStmt:
				return fmt.Sprintf(tipWFReturnErr, n.Value)
			case *parser.EventSentStmt:
				return fmt.Sprintf(tipWFEventSent, n.Value)
			*/
			default:

				return fmt.Sprintf("undokumentierte Literalverwendung: %T", token.Node)
			}
		}

		switch n.Parent().(type) {
		case *parser.FnStmtIf:
			return fmt.Sprintf(tipWFIf, n.Value)
			/*case *parser.WhileStmt:
			return fmt.Sprintf(tipWFWhile, n.Value)*/
		}

		return fmt.Sprintf("undokumentierte Identverwendung: %T", token.Node)

	default:
		return fmt.Sprintf("%T", token.Node)
	}

}

func shortStringLit(s string) string {
	const limit = 30
	s = strings.Split(s, ".")[0]
	s = strings.Split(s, "\n")[0]
	if len(s) < limit {
		return s
	}

	return s[:limit] + "..."
}
