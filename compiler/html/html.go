package html

import (
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/web/editor"
	"os"
	"path/filepath"
	"strings"
)

func Write(dstDir string, src *parser.Workspace) error {
	html := Render(src)
	outfile := dstDir
	if !strings.HasSuffix(outfile, ".html") {
		outfile = filepath.Join(outfile, "index.html")
	}

	return os.WriteFile(outfile, []byte(html), os.ModePerm)
}

func Render(src *parser.Workspace) string {
	var model editor.EditorPreview
	model.VSCode.ScriptUris = append(model.VSCode.ScriptUris, "https://cdn.tailwindcss.com")

	if src.Error != nil {
		model.Error = src.Error.Error()
	}

	lint := editor.Linter(func(ws *parser.Workspace) []linter.Hint {
		return linter.Lint(ws)
	})

	return editor.RenderViewHtml(lint, src, model)
}
