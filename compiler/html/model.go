package html

import (
	"github.com/worldiety/dddl/parser"
	"html/template"
	"time"
)

type Head struct {
	Nonce      string
	ScriptUris []string
}

type PreviewModel struct {
	DevMode              bool `json:"-"`
	Title                string
	Doc                  *Doc
	Hints                []template.HTML
	NamedTasks           []NamedTasks
	EditorText           string
	Error                string
	LastSaved            string
	Head                 Head
	LocalWorkspacePrefix string
	ProjectPlan          *ProjectPlan
}

type ProjectPlan struct {
	GanttChartSVG template.HTML
	Tasks         []*ProjectTask
}

type ProjectTask struct {
	Name     string
	Refs     []string
	Duration time.Duration
	Requires []string
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
	Name                string
	ShortDef            template.HTML
	Ref                 string
	Aggregates          []*Aggregate
	Types               []*Type
	Definition          template.HTML
	WorkPackageName     string
	WorkPackageRequires []string
	WorkPackageDuration string
}

func (c *Context) IsContext() bool {
	return true
}

func (c *Context) GroupTypesByCategory(cat string) []*Type {
	return FilterByCategory(c.Types, cat)
}

func FilterByCategory(types []*Type, cat string) []*Type {
	var res []*Type
	for _, t := range types {
		if t.Category == cat {
			res = append(res, t)
		}
	}

	return res
}

type Aggregate struct {
	Context             *Context `json:"-"`
	Category            string
	Name                string
	Ref                 string
	Types               []*Type
	Definition          template.HTML
	WorkPackageName     string
	WorkPackageRequires []string
	WorkPackageDuration string
}

func (c *Aggregate) GroupTypesByCategory(cat string) []*Type {
	return FilterByCategory(c.Types, cat)
}

func (c *Aggregate) IsContext() bool {
	return false
}

type Type struct {
	Node                parser.NamedType `json:"-"` // e.g. *parser.Struct, *parser.Choice, *parser.Function etc.
	Parent              any              `json:"-"` // either Context or Aggregate
	Category            string
	Name                string
	Ref                 string
	Definition          template.HTML
	SVG                 template.HTML
	SVGExtended         template.HTML
	Usages              []Usage
	WorkPackageName     string
	WorkPackageRequires []string
	WorkPackageDuration string
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
