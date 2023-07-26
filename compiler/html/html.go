package html

import (
	"embed"
	"github.com/worldiety/dddl/html"
	"github.com/worldiety/dddl/parser"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
)

//go:embed *.gohtml
var appFiles embed.FS

var viewFunc html.ViewFunc[PreviewModel]

func init() {
	viewFunc = html.MustParse[PreviewModel](
		html.FS(appFiles),
		html.Execute("ViewPage"),
	)
}

func RenderViewHtml(doc *parser.Workspace, model PreviewModel) string {
	model.Doc = transform(doc)
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
