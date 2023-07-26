package resolver

import (
	_ "embed"
	"github.com/worldiety/dddl/parser"
	"testing"
)

//go:embed testdata/test3.txt
var test string

func TestNewResolver(t *testing.T) {
	ws, err := parser.ParseWorkspaceText(map[string]string{"testdata/test3.txt": test})

	if err != nil {
		t.Fatal(err)
	}

	r := NewResolver(ws)

	if slice := Collect[*parser.Context](r); len(slice) != 4 {
		t.Fatal(slice)
	}

	if slice := Collect[*parser.Aggregate](r); len(slice) != 1 {
		t.Fatal(slice)
	}
}
