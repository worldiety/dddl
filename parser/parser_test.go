package parser

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"
)

//go:embed testdata/test3.txt
var test string

func TestParse(t *testing.T) {
	v, err := ParseText("testdata/test3.txt", test)
	if err != nil {
		t.Fatal(err)
	}

	buf, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}

func TestParse2(t *testing.T) {
	v, err := ParseWorkspaceText(map[string]string{"testdata/test3.txt": test})
	v.ResolveData(&Ident{
		Value: "asdf",
	})
	if err != nil {
		t.Fatal(err)
	}

	buf, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}
