package tpl

import (
	"bytes"
	"fmt"
	model2 "github.com/worldiety/dddl/compiler/model"
	"github.com/worldiety/dddl/parser"
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
			s = strings.TrimSpace(s)
			if s == "" {
				return ""
			}
			tmp := ""
			for _, line := range strings.Split(s, "\n") {
				tmp += "// " + line + "\n"
			}

			return tmp
		},

		"orTypeDefs": func(idents []*model2.TypeDef) string {
			tmp := ""
			for i, ident := range idents {
				tmp += typeDef(ident)
				if i < len(idents)-1 {
					tmp += " | "
				}
			}
			return tmp
		},

		"typeDef": typeDef,

		"typeDefAsIdent": typeDefAsIdent,
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

func typeDefAsIdent(q model2.QualifiedName) string {
	if q.IsUniverse() || q.Local {
		return model2.MakeUpIdentifier(q.Name)
	}

	if q.PackageName != "" {
		return model2.MakeUpIdentifier(q.PackageName + "." + q.Name)
	}

	return model2.MakeUpIdentifier(q.Name)
}

func typeDef(def *model2.TypeDef) string {
	if def.FuncDef != nil {
		tmp := "func("
		for _, t := range def.FuncDef.Input {
			tmp += typeDef(t)
			tmp += ","
		}
		tmp += ") ("

		tmp += typeDef(def.FuncDef.Output)
		tmp += ","
		tmp += typeDef(def.FuncDef.Error)
		tmp += ")"

		return tmp
	}

	tmp := typeDefName(def.Name)
	if tmp == "map" && len(def.Parameter) == 2 {
		return "map[" + typeDef(def.Parameter[0]) + "]" + typeDef(def.Parameter[1])
	}

	for i, t := range def.Parameter {
		tmp += typeDef(t)
		if i < len(def.Parameter)-1 {
			tmp += ","
		}

	}

	return tmp
}

func typeDefName(q model2.QualifiedName) string {
	if q.Local {
		return q.Name
	}

	if q.PackageName != "" {
		return q.PackageName + "." + q.Name
	}

	if q.IsUniverse() {
		switch q.Name {
		case parser.UString:
			return "string"
		case parser.UInt:
			return "int64"
		case parser.UList:
			return "[]"
		case parser.UMap:
			return "map"
		default:
			return "/*fix me*/ " + q.Name
		}
	}

	return q.Name
}
