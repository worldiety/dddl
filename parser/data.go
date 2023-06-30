package parser

// A Data is either a choice or a compound data type.
// Combining both is probably hard to understand.
// Without massive lookahead, we cannot distinguish that, so we will
// check that using a linter later.
type Data struct {
	node
	KeywordData *KeywordData       `@@`
	Name        *Ident             ` @@ ( "{"`
	Definition  *Definition        `@@?`
	ToDo        *ToDo              `@@? `
	First       *TypeDeclaration   ` (@@ `
	Fields      []*TypeDeclaration `("und" @@)*`
	Choices     []*TypeDeclaration `("oder" @@)*)?  "}")?`
}

func (d *Data) Empty() bool {
	return d.First == nil && len(d.Fields) == 0 && len(d.Choices) == 0
}

func (d *Data) Children() []Node {
	var res []Node
	res = append(res, d.KeywordData, d.Name)

	for _, declaration := range d.ChoiceTypes() {
		res = append(res, declaration)
	}

	for _, declaration := range d.FieldTypes() {
		res = append(res, declaration)
	}

	if d.Definition != nil {
		res = append(res, d.Definition)
	}

	if d.ToDo != nil {
		res = append(res, d.ToDo)
	}

	return res
}

// ChoiceTypes is nil, if any field is defined.
// If neither choices nor fields are defined, this also returns nil.
// This avoids visiting nodes twice.
func (d *Data) ChoiceTypes() []*TypeDeclaration {
	if len(d.Fields) > 0 || len(d.Choices) == 0 {
		return nil
	}

	var choices []*TypeDeclaration
	choices = append(choices, d.First)
	choices = append(choices, d.Choices...)
	return choices
}

// FieldTypes is nil if any choice type is defined, otherwise contains
// at least the First declaration.
func (d *Data) FieldTypes() []*TypeDeclaration {
	if len(d.Choices) > 0 {
		return nil
	}

	var fields []*TypeDeclaration
	if d.First != nil {
		fields = append(fields, d.First)
	}
	fields = append(fields, d.Fields...)
	return fields
}
