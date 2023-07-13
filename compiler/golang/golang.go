package golang

import (
	"fmt"
	gen "github.com/worldiety/dddl/compiler/golang/v1"
	"github.com/worldiety/dddl/parser"
)

func Write(opts Options, dstDir string, src *parser.Workspace) error {
	switch opts.Style {
	case DDDv1:
		if opts.VersionedGenOptions.DDDv1 == nil {
			gopt := gen.DefaultOptions(dstDir)
			opts.VersionedGenOptions.DDDv1 = &gopt
		}

		return gen.Write(*opts.VersionedGenOptions.DDDv1, dstDir, src)
	default:
		return fmt.Errorf("unsupported generator style: %s", opts.Style)
	}
}
