package editor

import (
	"github.com/worldiety/dddl/parser"
	"gitlab.worldiety.net/tschinke/html"
	"net/http"
	"net/http/httptest"
)

var viewFunc html.ViewFunc[EditorPreview]

func init() {
	viewFunc = html.MustParse[EditorPreview](
		html.FS(appFiles),
		html.Execute("ViewPage"),
	)
}

func RenderViewHtml(lint Linter, doc *parser.Doc, model EditorPreview) string {
	model = lintOnly(doc, lint, model)
	w := httptest.NewRecorder()

	viewFunc(w, &http.Request{}, model)
	return w.Body.String()
}
