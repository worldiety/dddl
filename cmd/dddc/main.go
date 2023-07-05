// Package main provides the executable for the ddd compiler.
package main

import (
	"flag"
	"fmt"
	"github.com/worldiety/dddl/compiler/html"
	"github.com/worldiety/dddl/parser"
	"log"
	"os"
)

const (
	HTML = "html"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	wd, _ := os.Getwd()
	format := flag.String("format", HTML, fmt.Sprintf("one of %s", HTML))
	in := flag.String("src", wd, "the source directory which contains the *.ddd files")
	out := flag.String("out", wd, "the target output directory or file")
	flag.Parse()

	ws, err := parser.ParseWorkspaceDir(*in)
	if ws == nil {
		return err
	}

	switch *format {
	case HTML:
		return html.Write(*out, ws)
	default:
		return fmt.Errorf("invalid format '%s'", *format)
	}
}
