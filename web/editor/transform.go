package editor

import (
	"bytes"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/puml"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	mdparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"golang.org/x/exp/slog"
	"html/template"
	"regexp"
	"strings"
)

func transform(pWS *parser.Workspace) *Doc {
	doc := &Doc{}

	for _, pdoc := range pWS.Docs() {
		for _, definition := range pdoc.Definitions {
			if definition.TypeDefinition != nil {
				if doc.SharedKernel == nil {
					doc.SharedKernel = &Context{}
				}

				if pdata := definition.TypeDefinition.DataType; pdata != nil {
					data := convertData(pdata)
					doc.SharedKernel.Data = append(doc.SharedKernel.Data, data)
				}

				if pWorkflow := definition.TypeDefinition.Workflow; pWorkflow != nil {
					wf := convertWorkflow(pWorkflow)
					doc.SharedKernel.Workflows = append(doc.SharedKernel.Workflows, wf)
				}
			}
		}
	}

	for _, pCtx := range pWS.CollectContextChildren() {
		ctx := &Context{}
		doc.Contexts = append(doc.Contexts, ctx)

		ctx.Name = pCtx.Name
		ctx.Definition = linkify(pCtx.Contexts[0], markdown(getContextDoc(pWS, ctx.Name)))
		ctx.Todo = linkify(pCtx.Contexts[0], markdown(getContextToDo(pWS, ctx.Name)))

		for _, child := range pCtx.Children {

			if pdata, ok := child.(*parser.Data); ok {
				data := convertData(pdata)
				ctx.Data = append(ctx.Data, data)
			}

			if pWorkflow, ok := child.(*parser.Workflow); ok {
				wf := convertWorkflow(pWorkflow)
				ctx.Workflows = append(ctx.Workflows, wf)
			}
		}

	}

	return doc
}

// getContextDoc joins all available definitions and todos into a big ball of text.
func getContextDoc(ws *parser.Workspace, name string) string {
	var sb strings.Builder
	for _, context := range ws.Contexts() {
		if context.Name.Value == name {
			sb.WriteString(getDoc(context))
		}
	}

	return sb.String()
}

func getDoc(def parser.Defineable) string {
	var sb strings.Builder
	if defText := parser.TextOf(def.GetDefinition()); defText != "" {
		sb.WriteString(defText)
		sb.WriteString("\n")
	}

	return sb.String()
}

// getContextDoc joins all available definitions and todos into a big ball of text.
func getContextToDo(ws *parser.Workspace, name string) string {
	var sb strings.Builder
	for _, context := range ws.Contexts() {
		if context.Name.Value == name {
			sb.WriteString(getTodo(context))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func getTodo(def parser.Defineable) string {
	var sb strings.Builder
	if defText := parser.TextOf(def.GetToDo()); defText != "" {
		sb.WriteString(defText)
	}

	return sb.String()
}

func convertData(pdata *parser.Data) *Data {
	data := &Data{}
	ws := parser.WorkspaceOf(pdata)
	data.Name = pdata.Name.Value
	q, _ := ws.Resolve(pdata.Name)
	data.Qualifier = q.String()

	if pdata.Definition != nil {
		data.Definition = linkify(pdata.Definition, markdown(pdata.Definition.Text))
	}

	if pdata.ToDo != nil {
		data.Todo = linkify(pdata.ToDo, markdown(pdata.ToDo.Text.Text))
	}

	if !pdata.Empty() {
		svg, err := plantuml.RenderLocal("svg", puml.Data(pdata))
		if err != nil {
			slog.Error("failed to convert data to puml", slog.Any("err", err))
		}

		data.SVG = template.HTML(svg)
	}

	return data
}

func convertWorkflow(pWorkflow *parser.Workflow) *Workflow {
	wf := &Workflow{}

	ws := parser.WorkspaceOf(pWorkflow)
	wf.Name = pWorkflow.Name.Value
	q, _ := ws.Resolve(pWorkflow.Name)
	wf.Qualifier = q.String()

	if pWorkflow.Definition != nil {
		wf.Definition = linkify(pWorkflow.Definition, markdown(pWorkflow.Definition.Text))
	}

	if pWorkflow.ToDo != nil {
		wf.Todo = linkify(pWorkflow.ToDo, markdown(pWorkflow.ToDo.Text.Text))
	}

	if pWorkflow.Block != nil && len(pWorkflow.Block.Statements) > 0 {
		svg, err := plantuml.RenderLocal("svg", puml.Workflow(pWorkflow))
		if err != nil {
			slog.Error("failed to convert workflow to puml", slog.Any("err", err))
		}

		wf.SVG = template.HTML(svg)
	}

	return wf
}

var regexWord = regexp.MustCompile(`([À-ž]|\w)+`)

func linkify(nearest parser.Node, text template.HTML) template.HTML {
	ws := parser.WorkspaceOf(nearest)
	tmp := regexWord.ReplaceAllStringFunc(string(text), func(s string) string {
		potentialIdent := parser.NewIdentWithParent(nearest, s)
		_, ok := ws.Resolve(potentialIdent)
		if !ok {
			return s
		}

		return href(ws, potentialIdent)
	})

	return template.HTML(tmp)
}

func markdown(text string) template.HTML {
	text = parser.TextOf(text)
	if text == "" {
		return ""
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			mdparser.WithAutoHeadingID(),
			mdparser.WithASTTransformers(util.PrioritizedValue{Value: &tailwindTransformer{}}),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		panic(err)
	}

	return template.HTML(buf.Bytes())
}

type tailwindTransformer struct {
}

func (t *tailwindTransformer) Transform(node *ast.Document, reader text.Reader, pc mdparser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if node, ok := n.(*ast.ListItem); ok {
			//node.SetAttributeString("class", []byte("list-disc"))
			_ = node
		}

		if node, ok := n.(*ast.List); ok {
			node.SetAttributeString("class", []byte("list-disc ml-8"))
		}

		if node, ok := n.(*ast.Paragraph); ok {
			node.SetAttributeString("class", []byte("mb-2 mt-2"))
		}

		if node, ok := n.(*ast.Heading); ok {
			node.SetAttributeString("class", []byte("font-medium"))
		}

		return ast.WalkContinue, nil
	})
}
