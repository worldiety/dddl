package editor

import (
	"embed"
	"fmt"
	"github.com/worldiety/dddl/html"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/slog"
	"html/template"
	"net/http"
	"time"
)

//go:embed *.gohtml
var appFiles embed.FS

type Saver func(text string) error

type Loader func() string

type Parser func(text string) (*parser.Workspace, error)

type Linter func(doc *parser.Workspace) []linter.Hint

func Handler(devMode bool, load Loader, save Saver, parse Parser, lint Linter) http.HandlerFunc {
	return html.Handler(
		html.MustParse[EditorPreview](
			html.FS(appFiles),
			html.Execute("EditorPage"),
		),
		html.OnRequest(
			func(r *http.Request, model EditorPreview) EditorPreview {
				model.devMode = devMode
				model.Title = "wdy visual workflow"
				model.EditorText = load()
				model = parseAndLint(parse, lint, model, model.EditorText)
				return model
			},
		),
		html.Update(
			html.CaseWithAlias("save", func(model EditorPreview, msg Save) EditorPreview {
				model.EditorText = msg.Text
				model.Error = ""
				model.Hints = nil
				model.LastSaved = "!!!"
				if err := save(msg.Text); err != nil {
					model.Error = err.Error()
					slog.Error("could not save", slog.Any("err", err))
					return model
				}

				model.LastSaved = time.Now().Format(time.DateTime)
				model = parseAndLint(parse, lint, model, msg.Text)
				return model
			}),
		),
	)
}

func parseAndLint(parse Parser, lint Linter, model EditorPreview, text string) EditorPreview {
	doc, err := parse(text)
	if err != nil {
		model.Error = err.Error()
		slog.Error("could not save", slog.Any("err", err))
		return model
	}

	model = lintOnly(doc, lint, model)

	return model
}

func lintOnly(doc *parser.Workspace, lint Linter, model EditorPreview) EditorPreview {
	model.Doc = transform(doc)
	//fmt.Printf("%#v", model)

	for _, hint := range lint(doc) {
		str := hint.String(func(ident *parser.Ident) string {
			return fmt.Sprintf(`<a class="text-green-600" href="#%s">%s</a>`, ident.Value, ident.Value)
		})
		model.Hints = append(model.Hints, template.HTML(str))
	}

	return model
}
