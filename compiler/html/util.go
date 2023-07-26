package html

import (
	"bytes"
	"github.com/worldiety/dddl/parser"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	mdparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"html/template"
)

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
