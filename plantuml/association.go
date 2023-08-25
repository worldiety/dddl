package plantuml

import "io"

type AssocType string

const (
	AssocAggregation = "o"
	AssocComposition = "*"
	AssocExtension   = "<|"
)

type Association struct {
	Owner            string
	OwnerCardinality string
	Child            string
	Type             AssocType
}

func (a *Association) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Print(`"`)
	w.Print(escapeP(a.Owner))
	w.Print(`"`)
	if a.OwnerCardinality != "" {
		w.Print(` "`)
		w.Print(a.OwnerCardinality)
		w.Print(`" `)
	}
	w.Print(string(a.Type))
	w.Print("-- ")
	w.Print(escapeP(a.Child))
	w.Print("\n")

	return w.Err
}
