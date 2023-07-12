package parser

// A Data is either a choice or a compound data type.
// Combining both is probably hard to understand.
// Without massive lookahead, we cannot distinguish that, so we will
// check that using a linter later.
type Data struct {
	node
	KeywordData *KeywordData `@@`
	Name        *Ident       ` @@ ( "{"`
	ToDo        *ToDo        `@@? `
	Definition  *Definition  `@@?`
	First       *TypeDef     ` (@@ `
	Fields      []*TypeDef   `("und" @@)*`
	Choices     []*TypeDef   `("oder" @@)*)?  "}")?`
}

func DataOf(root Node) *Data {
	for root != nil {
		if d, ok := root.(*Data); ok {
			return d
		}
		root = root.Parent()
	}

	return nil
}

func (d *Data) DataOrWorkflow() bool {
	return true
}

func (d *Data) GetDefinition() string {
	return d.Definition.Value()
}

func (d *Data) GetToDo() string {
	return d.ToDo.Value()
}

func (d *Data) DeclaredName() *Ident {
	return d.Name
}

func (d *Data) Qualifier() Qualifier {
	return Qualifier{
		Context: d.Parent().(*Context),
		Name:    d.Name,
	}
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

func (d *Data) IsChoiceType() bool {
	return len(d.ChoiceTypes()) > 0
}

// ChoiceTypes is nil, if any field is defined.
// If neither choices nor fields are defined, this also returns nil.
// This avoids visiting nodes twice.
func (d *Data) ChoiceTypes() []*TypeDef {
	if len(d.Fields) > 0 || len(d.Choices) == 0 {
		return nil
	}

	var choices []*TypeDef
	choices = append(choices, d.First)
	choices = append(choices, d.Choices...)
	return choices
}

// FieldTypes is nil if any choice type is defined, otherwise contains
// at least the First declaration.
func (d *Data) FieldTypes() []*TypeDef {
	if len(d.Choices) > 0 {
		return nil
	}

	var fields []*TypeDef
	if d.First != nil {
		fields = append(fields, d.First)
	}
	fields = append(fields, d.Fields...)
	return fields
}
