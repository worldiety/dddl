package html

import (
	"embed"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
	"github.com/worldiety/hg"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

//go:embed *.gohtml
var appFiles embed.FS

var viewFunc hg.ViewFunc[PreviewModel]

func init() {
	viewFunc = hg.MustParse[PreviewModel](
		hg.FS(appFiles),
		hg.Execute("ViewPage"),
	)

}

func RenderViewHtml(ws *parser.Workspace, model PreviewModel) string {
	rslv := resolver.NewResolver(ws)
	model.Doc = transform(rslv, model)
	lintHints := linter.Lint(rslv)
	model = transformLintHints(rslv, lintHints, model)
	model.ProjectPlan = newProjectPlan(rslv, model)

	w := httptest.NewRecorder()

	viewFunc(w, &http.Request{}, model)
	return w.Body.String()
}

func Write(dstDir string, src *parser.Workspace) error {
	buf := Render(src)
	outfile := dstDir
	if !strings.HasSuffix(outfile, ".html") {
		outfile = filepath.Join(outfile, "index.html")
	}

	return os.WriteFile(outfile, []byte(buf), os.ModePerm)
}

func Render(src *parser.Workspace) string {
	var model PreviewModel
	model.Head.ScriptUris = append(model.Head.ScriptUris, "https://cdn.tailwindcss.com")

	if src.Error != nil {
		model.Error = src.Error.Error()
	}

	return RenderViewHtml(src, model)
}
