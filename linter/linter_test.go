package linter

import (
	_ "embed"
	"fmt"
	"github.com/worldiety/dddl/parser"
	"testing"
)

//go:embed testdata/test3.txt
var test string

func TestLint(t *testing.T) {
	v, err := parser.ParseText("testdata/test3.txt", test)
	if err != nil {
		t.Fatal(err)
	}

	ws := &parser.Workspace{Documents: map[string]*parser.Doc{"bla": v}}

	hints := Lint(ws)
	fmt.Printf("%#v", hints)
}
