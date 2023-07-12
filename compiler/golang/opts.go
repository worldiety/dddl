package golang

import v1gen "github.com/worldiety/dddl/compiler/golang/v1"

type GenStyle string

const (
	DDDv1 GenStyle = "dddv1"
)

type VersionedGenOptions struct {
	DDDv1 *v1gen.Options `json:"dddv1"`
}

type Options struct {
	Style               GenStyle            `json:"style"`
	VersionedGenOptions VersionedGenOptions `json:"options"`
}

func Default() Options {
	return Options{
		Style: DDDv1,
	}
}
