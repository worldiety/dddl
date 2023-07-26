package plantuml

import "io"

const ThemeCerulean = "https://raw.githubusercontent.com/bschwarz/puml-themes/master/themes/cerulean/puml-theme-cerulean.puml"

type DefaultTextAlignment string

const (
	DTACenter = "skinparam defaulttextalignment center"
)

type Diagram struct {
	includes             []string
	Renderables          []Renderable
	DefaultTextAlignment DefaultTextAlignment
	BackgroundColor      string
}

func NewDiagram() *Diagram {
	d := &Diagram{}
	return d
}

func (d *Diagram) Add(r ...Renderable) *Diagram {
	d.Renderables = append(d.Renderables, r...)
	return d
}

func (d *Diagram) Include(inc ...string) *Diagram {
	d.includes = append(d.includes, inc...)
	return d
}

func (d *Diagram) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Print("@startuml\n")
	if d.BackgroundColor != "" {
		w.Print("skinparam backgroundColor ")
		w.Print(d.BackgroundColor)
		w.Print("\n")
	}

	if d.DefaultTextAlignment != "" {
		w.Print(string(d.DefaultTextAlignment))
		w.Print("\n")
	}

	for _, include := range d.includes {
		w.Print("!include ")
		w.Print(include)
		w.Print("\n")
	}

	for _, renderable := range d.Renderables {
		if err := renderable.Render(wr); err != nil {
			return err
		}
	}

	w.Print("@enduml\n")
	return w.Err
}
