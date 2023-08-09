package html

import (
	"embed"
	"errors"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/hg"
	"html/template"
	"strings"
)

//go:embed erroroverlay.gohtml
var erroroverlay embed.FS
var errorOverlayFunc hg.ViewFunc[ErrorModel]

type ErrorModel struct {
	Messages []template.HTML
}

func init() {
	errorOverlayFunc = hg.MustParse[ErrorModel](
		hg.FS(erroroverlay),
		hg.Execute("erroroverlay"),
	)

}

// PostInsertError looks wrong (and maybe is), however live preview renderers would like to keep the last
// valid rendering and place an error screen on-top, so that the scroll-position etc. will not change while
// typing.
func PostInsertError(html string, err error) string {
	model := ErrorModel{}

	var dperr *parser.DocParserError
	if errors.As(err, &dperr) {
		for file, err := range dperr.Errors {
			file = strings.TrimPrefix(file, "file://")
			model.Messages = append(model.Messages, template.HTML(`<p class="text-rose-700">`+file+`</p><p class="pl-4 text-rose-700">`+err.Error()+"</p>"))
		}
	} else {
		model.Messages = append(model.Messages, template.HTML(dperr.Error()))
	}

	dlg := errorOverlayFunc.ToString(model)
	return strings.Replace(html, `<div id="errorOverlay"></div>`, dlg, 1)
}
