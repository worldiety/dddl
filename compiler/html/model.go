package html

import (
	"html/template"
)

type Head struct {
	Nonce      string
	ScriptUris []string
}

type PreviewModel struct {
	DevMode    bool `json:"-"`
	Title      string
	Doc        *Doc
	Hints      []template.HTML
	NamedTasks []NamedTasks
	EditorText string
	Error      string
	LastSaved  string
	Head       Head
}

type NamedTasks struct {
	Name  string
	Tasks []template.HTML
}

type Doc struct {
	SharedKernel *Context
	Contexts     []*Context
}

type Context struct {
	Name       string
	ShortDef   template.HTML
	Ref        string
	Aggregates []*Aggregate
	Types      []*Type
	Definition template.HTML
}

type Aggregate struct {
	Context    *Context `json:"-"`
	Category   string
	Name       string
	Ref        string
	Types      []*Type
	Definition template.HTML
}

type Type struct {
	Aggregate  *Aggregate `json:"-"`
	Context    *Context   `json:"-"`
	Category   string
	Name       string
	Ref        string
	Definition template.HTML
	SVG        template.HTML
	Usages     []Usage
}

type Usage struct {
	Name string
	Ref  string
}

type Workflow struct {
	Name       string
	Qualifier  string
	Definition template.HTML
	Todo       template.HTML
	Choices    []string
	SVG        template.HTML
}
