// Package main provides the executable dddw for the standalone ddd web ui.
package main

import (
	"flag"
	"fmt"
	"github.com/worldiety/dddl/linter"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/web"
	"github.com/worldiety/dddl/web/editor"
	"golang.org/x/exp/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
	"unicode/utf8"
)

func main() {
	host := flag.String("host", "localhost:8080", "the host and port to bind on")
	dev := flag.Bool("devMode", false, "set only to true, if your are a developer, e.g. for hot-reload from IDE")
	fname := flag.String("file", "Modell.txt", "file to save and load the model")
	flag.Parse()

	saver := editor.Saver(func(text string) error {
		shortName := filepath.Base(*fname)
		historyDir := "." + shortName + ".versions"
		if err := os.MkdirAll(historyDir, os.ModePerm); err != nil {
			slog.Error("failed to create version dir", slog.Any("err", err))
		}

		now := time.Now()
		historySnapshot := filepath.Join(historyDir, now.Format(time.RFC3339)+"_"+shortName)
		if err := os.WriteFile(historySnapshot, []byte(text), os.ModePerm); err != nil {
			slog.Error("failed to write history entry", slog.Any("err", err))
		}

		tmpf := *fname + "_" + strconv.Itoa(int(now.UnixMilli())) + ".tmp"
		if err := os.WriteFile(tmpf, []byte(text), os.ModePerm); err != nil {
			return err
		}

		if err := os.Rename(tmpf, *fname); err != nil {
			return err
		}

		return nil
	})

	parse := editor.Parser(func(text string) (*parser.Workspace, error) {
		fname := "???"
		doc, err := parser.ParseText(fname, text)
		if err != nil {
			return nil, fmt.Errorf("cannot parse model: %w", err)
		}

		return &parser.Workspace{Documents: map[string]*parser.Doc{fname: doc}}, nil
	})

	linter := editor.Linter(func(doc *parser.Workspace) []linter.Hint {
		return linter.Lint(doc)
	})

	loader := editor.Loader(func() string {
		var text string
		buf, err := os.ReadFile(*fname)
		if err == nil {
			if utf8.Valid(buf) {
				text = string(buf)
			}
		}

		return text
	})

	if !*dev {
		go func() {
			time.Sleep(500 * time.Millisecond)
			uglyLaunchMacOSChromeBrowserWindow(*host)
		}()
	}

	if err := web.StartServer(*host, *dev, loader, saver, parse, linter); err != nil {
		panic(err)
	}
}

func uglyLaunchMacOSChromeBrowserWindow(host string) {
	///Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome --new-window --app=https://localhost:8080
	cmd := exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--new-window", "--app=http://"+host)
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}
}
