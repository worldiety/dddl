// Package main provides the executable for the ddd compiler.
package main

import (
	"flag"
	"fmt"
	"github.com/worldiety/dddl/compiler/golang"
	"github.com/worldiety/dddl/compiler/html"
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/slog"
	"log"
	"os"
)

const (
	HTML    = "html"
	Grammar = "grammar"
	Go      = "go"
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

	if ws.Error != nil {
		slog.Error("cannot parse workspace *.ddd files properly", slog.Any("err", ws.Error))
	}

	switch *format {
	case HTML:
		return html.Write(*out, ws)
	case Grammar:
		fmt.Println(parser.NewParser().String())
		return nil
	case Go:
		opts := golang.Default()
		return golang.Write(opts, *out, ws)
	default:
		return fmt.Errorf("invalid format '%s'", *format)
	}
}
