package parser

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	v, err := Parse("testdata/test3.txt")
	if err != nil {
		t.Fatal(err)
	}

	buf, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf))
}
