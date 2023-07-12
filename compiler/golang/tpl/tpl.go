package tpl

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"strings"
	"text/template"
)

func Format(s string) string {
	buf, err := format.Source([]byte(s))
	if err != nil {
		panic(fmt.Sprintf("%s\n%v", buf, err))
	}

	return string(buf)
}

func Execute(fsys fs.FS, name string, model any) ([]byte, error) {
	tpl := template.New(name)

	tpl.Funcs(map[string]any{
		"makeComment": func(s string) string {
			tmp := ""
			for _, line := range strings.Split(string(s), "\n") {
				tmp += "// " + line + "\n"
			}

			return tmp
		},

		"orIdents": func(idents []string) string {
			tmp := ""
			for i, ident := range idents {
				tmp += string(ident)
				if i < len(idents)-1 {
					tmp += " | "
				}
			}
			return tmp
		},
	})

	tpl, err := tpl.ParseFS(fsys, "*.gotmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, name, model); err != nil {
		return nil, fmt.Errorf("cannot execute template: %w", err)
	}

	b, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(buf.String())
		return nil, fmt.Errorf("cannot format src: %w", err)
	}

	return b, nil
}
