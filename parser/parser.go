package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"golang.org/x/exp/slices"
	"os"
	"strings"
)

// TODO manuelle TODOS
type Doc struct {
	node
	Contexts []*Context `@@*`
}

func (n *Doc) DataByName(name string) *Data {
	for _, context := range n.Contexts {
		for _, element := range context.Elements {
			if element.DataType != nil {
				if element.DataType.Name.Name == name {
					return element.DataType
				}
			}
		}
	}

	return nil
}

func (n *Doc) WorkflowByName(name string) *Workflow {
	for _, context := range n.Contexts {
		for _, element := range context.Elements {
			if element.Workflow != nil {
				if element.Workflow.Name.Name == name {
					return element.Workflow
				}
			}
		}
	}

	return nil
}

func (n *Doc) Children() []Node {
	var res []Node
	for _, context := range n.Contexts {
		res = append(res, context)
	}

	return res
}

type Context struct {
	node
	KeywordContext *KeywordContext   `@@`
	Name           *Ident            `@@`
	ToDo           *ToDo             `@@?`
	Definition     *Definition       `@@?`
	Elements       []*TypeDefinition `@@*`
}

func (n *Context) DataTypes() []*Data {
	var res []*Data
	for _, element := range n.Elements {
		if element.DataType != nil {
			res = append(res, element.DataType)
		}
	}

	slices.SortFunc(res, func(a, b *Data) bool {
		return a.Name.Name < b.Name.Name
	})

	return res
}

func (n *Context) Workflows() []*Workflow {
	var res []*Workflow
	for _, element := range n.Elements {
		if element.Workflow != nil {
			res = append(res, element.Workflow)
		}
	}

	slices.SortFunc(res, func(a, b *Workflow) bool {
		return a.Name.Name < b.Name.Name
	})

	return res
}

func (n *Context) Children() []Node {
	var res []Node
	res = append(res, n.KeywordContext, n.Name)
	if n.ToDo != nil {
		res = append(res, n.ToDo)
	}

	for _, element := range n.Elements {
		res = append(res, element)
	}

	if n.Definition != nil {
		res = append(res, n.Definition)
	}

	return res
}

// A TypeDefinition is either a DataType or a Workflow.
// To simplify the parsing without lookahead, we just use
// this kind of union.
// See also TypeDeclaration.
type TypeDefinition struct {
	node

	DataType *Data     `@@`
	Workflow *Workflow `|@@`
}

func (d *TypeDefinition) Children() []Node {
	if d.DataType != nil {
		return []Node{d.DataType}
	}

	if d.Workflow != nil {
		return []Node{d.Workflow}
	}

	return nil
}

func (d *TypeDefinition) Name() *Ident {
	if d.DataType != nil {
		return d.DataType.Name
	}

	return d.Workflow.Name
}

// A TypeDeclaration is either a Name or a parameterized Name
// - a generic.
type TypeDeclaration struct {
	node

	Name   *Ident             `@@`
	Params []*TypeDeclaration `("[" @@ ("," @@)* "]" )?`
}

func (n *TypeDeclaration) Children() []Node {
	if n == nil {
		return nil
	}

	var res []Node
	res = append(res, n.Name)
	for _, param := range n.Params {
		res = append(res, param)
	}

	return res
}

// TextOf extracts and normalizes string literals.
func TextOf(s string) string {
	s = strings.Trim(s, " \n\t")

	var tmp []string
	for _, line := range strings.Split(s, "\n") {
		tmp = append(tmp, strings.TrimSpace(line))
	}

	return strings.Join(tmp, "\n")
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
	return v, nil
}

func ParseText(text string) (*Doc, error) {
	parser := NewParser()
	return parser.ParseBytes("Modell.wvw", []byte(text))
}

func NewParser() *participle.Parser[Doc] {
	var basicLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"comment", `//.*|/\*.*?\*/`},
		{"Text", `\"(\\.|[^"\\])*\"`},
		{"Name", `([À-ž]|\w)+`},
		{"Assign", `=`},
		{"Colon", "[:,]"},
		{"Block", "[{}]"},
		{"Generic", `[\[\]\(\)]`},
		{"whitespace", `[ \t\n\r]+`},
	})

	parser, err := participle.Build[Doc](
		participle.Lexer(basicLexer),
		participle.Unquote("Text"),
	)

	if err != nil {
		panic(err) // this is always a programming error
	}

	return parser
}
