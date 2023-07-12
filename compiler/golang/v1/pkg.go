package gen

import (
	"github.com/worldiety/dddl/compiler/golang/tpl"
	"github.com/worldiety/dddl/compiler/model"
)

func (g *gen) renderPkg(pkg *model.Package) ([]byte, error) {
	return tpl.Execute(templates, "Package", pkg)
}
