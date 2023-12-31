package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TextOf extracts and normalizes string literals.
func TextOf(s string) string {
	lines := strings.Split(s, "\n")
	var tmp []string
	indent := -1
	for _, line := range lines {
		trimLine := strings.TrimSpace(line)
		if trimLine == "" && len(tmp) == 0 {
			continue
		}

		if trimLine == "" {
			tmp = append(tmp, "")
			continue
		}

		if indent == -1 {
			after := strings.TrimLeft(line, " ")
			indent = len(line) - len(after)
		}

		tmp = append(tmp, negativeIndent(indent, line))
	}

	return strings.Join(tmp, "\n")
}

func negativeIndent(indent int, s string) string {
	var sb strings.Builder
	for i, r := range s {
		if i < indent && r == ' ' {
			continue
		}

		sb.WriteRune(r)
	}

	return sb.String()
}

func Parse(fname string) (*Doc, error) {

	buf, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	parser := NewParser()
	v, err := parser.ParseBytes(fname, buf)
	if err != nil {
		return nil, err
	}

	//fmt.Println(parser.String())
	attachParent(nil, v)
	return v, nil
}

func ParseText(filename, text string) (*Doc, error) {
	parser := NewParser()
	doc, err := parser.ParseBytes(filename, []byte(text))
	if doc != nil {
		attachParent(nil, doc)
	}

	return doc, err
}

func ParseWorkspaceDir(dir string) (*Workspace, error) {
	files := map[string]string{}
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if strings.HasSuffix(path, ".ddd") {
			buf, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read '%s': %w", path, err)
			}

			files[path] = string(buf)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("cannot walk '%s': %w", dir, err)
	}

	return ParseWorkspaceText(files)
}

// ParseWorkspaceText loads from filename->text and tries to parse each one.
// Continues and always returns a Workspace, even if error is not nil.
// If error is not nil, it is always [DocParserError].
func ParseWorkspaceText(filenamesWithText map[string]string) (*Workspace, error) {
	var tmp *DocParserError
	parserErr := func() *DocParserError {
		if tmp == nil {
			tmp = &DocParserError{Errors: map[string]error{}}
		}

		return tmp
	}

	workspace := &Workspace{Documents: map[string]*Doc{}}
	filenames := maps.Keys(filenamesWithText)
	slices.Sort(filenames)
	for _, filename := range filenames {
		doc, err := ParseText(filename, filenamesWithText[filename])
		if err != nil {
			parserErr().Errors[filename] = err
			continue
		}

		workspace.Documents[filename] = doc
	}

	attachParent(nil, workspace)
	if tmp != nil {
		workspace.Error = tmp
	}

	if tmp != nil {
		return workspace, tmp
	}

	return workspace, nil
}

type DocParserError struct {
	Errors map[string]error
}

func (e *DocParserError) Error() string {
	tmp := "DocParserError"
	keys := maps.Keys(e.Errors)
	slices.Sort(keys)
	for _, key := range keys {
		tmp += fmt.Sprintf(" * failed '%s': %s\n", key, e.Errors[key])
	}

	return tmp
}

func (e *DocParserError) Unwrap() []error {
	return maps.Values(e.Errors)
}

func NewParser() *participle.Parser[Doc] {
	var basicLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"comment", `//.*|/\*.*?\*/`},
		{"Keyword", `(?i)\b(Auswahl|choice|Daten|data|Synonym|alias|type|Typ|task|Aufgabe|solange|while)`},
		{"Text", `\"(\\.|[^"\\])*\"`},
		{"Name", `([À-ž]|\w)+`},
		{"Assign", `=`},
		{"FnRet", "->"},
		{"Option", `\?`},
		{"Colon", `[:,><.|@]`},
		{"Block", "[{}]"},
		{"Generic", `[\[\]\(\)]`},
		{"whitespace", `[ \t\n\r]+`},
	})

	parser, err := participle.Build[Doc](
		participle.Lexer(basicLexer),
		participle.Unquote("Text"),
		participle.Union[FnStmt](&FnStmtIf{}, &FnStmtBlock{}, &FuncTypeRet{}, &FnLitExpr{}, &FnStmtWhile{}),
		participle.Union[NamedType](&Function{}, &Choice{}, &Struct{}, &Type{}, &Alias{}, &Context{}, &Aggregate{}),
	)

	if err != nil {
		panic(err) // this is always a programming error
	}

	return parser
}

func attachParent(parent, n Node) {
	if n == nil {
		return
	}

	n.setParent(parent)
	for _, c := range n.Children() {
		attachParent(n, c)
	}
}
