package editor

import (
	"fmt"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/parser"
	"testing"
)

func TestRenderViewHtml(t *testing.T) {
	doc, err := parser.ParseText("test", `Kontext HelloWorld`)
	if err != nil {
		t.Fatal(err)
	}

	str := RenderViewHtml(func(doc *parser.Doc) []linter.Hint {
		return nil
	}, doc, EditorPreview{})
	fmt.Println(str)
}
