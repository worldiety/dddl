package asciidoc

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/resolver"
	"strings"
)

func Render(src *parser.Workspace) string {
	rslv := resolver.NewResolver(src)
	var out strings.Builder
	out.WriteString("= Implement me\n\n")
	for _, context := range rslv.Contexts() {
		out.WriteString("== ")
		out.WriteString(context.Name)
		out.WriteString("\n")
	}

	return out.String()
}
