package gen

import (
	"embed"
	"fmt"
	"github.com/worldiety/dddl/compiler/golang/tpl"
	"github.com/worldiety/dddl/compiler/model"
	"github.com/worldiety/dddl/parser"
	"log"
	"path/filepath"
	"strings"
	"testing/fstest"
)

//go:embed *.gotmpl
var templates embed.FS

func Write(opts Options, dstDir string, src *parser.Workspace) error {
	g := &gen{
		opts:   opts,
		dstDir: dstDir,
		src:    src,
		fs:     map[string]*fstest.MapFile{},
	}

	g.genAll()
	if len(g.errors) == 0 {
		if err := g.emit(); err != nil {
			return fmt.Errorf("cannot generate src: %w", err)
		}
	}

	if len(g.errors) > 0 {
		return g.errors[0]
	}

	log.Println(g.String())

	// intentionally we use opts.Dir and not dstDir because we may have an awkward location
	// and auto-detect by default the module root and place the files into a stable position.
	if err := model.Deploy(opts.Dir, g.fs); err != nil {
		return fmt.Errorf("cannot deploy generated files: %w", err)
	}

	return nil
}

type gen struct {
	opts   Options
	dstDir string
	src    *parser.Workspace
	fs     fstest.MapFS
	pkgs   []*model.Package
	errors []error
}

func (g *gen) String() string {
	var sb strings.Builder
	for name, file := range g.fs {
		sb.WriteString(name)
		sb.WriteString(":\n\n")
		sb.WriteString(string(file.Data))
		sb.WriteString("\n\n\n")
	}

	return sb.String()
}

func (g *gen) exec(fname string, tplname string, model any) {
	buf, err := tpl.Execute(templates, tplname, model)
	if err != nil {
		g.errors = append(g.errors, err)
	}

	g.fs[fname] = &fstest.MapFile{
		Data: buf,
	}
}

func (g *gen) emit() error {
	for _, pkg := range g.pkgs {
		buf, err := g.renderPkg(pkg)
		if err != nil {
			return fmt.Errorf("cannot render package %s: %w", pkg.Name, err)
		}

		fname := filepath.Join(g.opts.DomainPackagePrefix, pkg.Name, "domain.gen.go")
		g.fs[fname] = &fstest.MapFile{
			Data: buf,
		}
	}
	return nil
}

func (g *gen) genAll() {
	g.pkgs = model.Convert(g.src)

}
