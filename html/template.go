package html

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slog"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

type htmlParser struct {
	tpls        *template.Template
	tplExecName string
	fsyss       []fs.FS
	funcs       template.FuncMap
}

type TplOption func(p *htmlParser)

func FS(fsys fs.FS) TplOption {
	return func(p *htmlParser) {
		p.fsyss = append(p.fsyss, fsys)
	}
}

// Execute will invoke the given template name instead of evaluating the entire template set anonymously.
func Execute(tplName string) TplOption {
	return func(p *htmlParser) {
		p.tplExecName = tplName
	}
}

func NamedFunc(name string, fun any) TplOption {
	return func(p *htmlParser) {
		p.funcs[name] = fun
	}
}

func MustParse[Model any](options ...TplOption) ViewFunc[Model] {
	f, err := Parse[Model](options...)
	if err != nil {
		panic(err)
	}

	return f
}

func Parse[Model any](options ...TplOption) (ViewFunc[Model], error) {
	parser := &htmlParser{
		tpls: template.New(""),
		funcs: template.FuncMap{
			// TODO still required?
			"toJSON": func(obj any) string {
				buf, err := json.Marshal(obj)
				if err != nil {
					panic(err)
				}
				return string(buf)
			},

			"map": func(args ...any) map[string]any {
				res := map[string]any{}
				for i := 0; i < len(args); i += 2 {
					res[args[i].(string)] = args[i+1]
				}

				return res
			},

			// unsafe func but required for writing custom inline "slot-definitions", see also str
			"html": func(args ...any) template.HTML {
				var sb strings.Builder
				for _, arg := range args {
					sb.WriteString(fmt.Sprintf("%v", arg))
				}

				return template.HTML(sb.String())
			},

			"str": func(args ...any) string {
				var sb strings.Builder
				for _, arg := range args {
					sb.WriteString(fmt.Sprintf("%v", arg))
				}

				return sb.String()
			},
		},
	}

	parser.funcs["evaluate"] = func(templateName string, data any) template.HTML {
		var tmp bytes.Buffer
		err := parser.tpls.ExecuteTemplate(&tmp, templateName, data)
		if err != nil {
			return template.HTML("<p>" + err.Error() + "</p>")
		}

		return template.HTML(tmp.String())
	}

	for _, option := range options {
		option(parser)
	}

	parser.tpls.Funcs(parser.funcs)

	for _, fsys := range parser.fsyss {
		_, err := parser.tpls.ParseFS(fsys, "*.gohtml")
		if err != nil {
			return nil, fmt.Errorf("cannot parse *.gohtml files: %w", err)
		}

	}

	return func(w http.ResponseWriter, r *http.Request, model Model) {
		if parser.tplExecName == "" {
			if err := parser.tpls.Execute(w, model); err != nil {
				slog.Error("cannot execute anonymous template", err)
				w.Write([]byte(fmt.Sprintf(`<p style="color:red">cannot parse anonymous template: %s</p>'`, err.Error())))
			}

			return
		}

		if err := parser.tpls.ExecuteTemplate(w, parser.tplExecName, model); err != nil {
			slog.Error("cannot execute named template", err)
			w.Write([]byte(fmt.Sprintf(`<p style="color:red">cannot parse '%s' template: %s</p>'`, parser.tplExecName, err.Error())))
		}
	}, nil

}
