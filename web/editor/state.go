package editor

import "html/template"

type VSCode struct {
	Nonce      string
	ScriptUris []string
}

type EditorPreview struct {
	devMode    bool
	Title      string
	Doc        *Doc
	Hints      []template.HTML
	EditorText string
	Error      string
	LastSaved  string
	VSCode     VSCode
}

func (m EditorPreview) DevMode() bool {
	return m.devMode
}

type Doc struct {
	Contexts []*Context
}

type Context struct {
	Name       string
	Definition template.HTML
	Todo       template.HTML
	Data       []*Data
	Workflows  []*Workflow
}

type Data struct {
	Name       string
	Definition template.HTML
	Todo       template.HTML
	Fields     []string
	SVG        template.HTML
}

type Workflow struct {
	Name       string
	Definition template.HTML
	Todo       template.HTML
	Choices    []string
	SVG        template.HTML
}
