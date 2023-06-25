package parser

type Workflow struct {
	node

	Name         *Ident      `"Arbeitsablauf" @@ ("=" `
	ToDo         *ToDo       `@@? `
	Dependencies *Input      `("Abhängigkeiten" ":" @@)?`
	Input        *Input      ` ("Eingabe" ":" @@`
	Output       *Output     `"Ausgabe" ":" @@`
	Block        *Stmts      `("Ablauf" "{" @@ "}")?)?`
	Definition   *Definition `@@?)?`
}

func (n *Workflow) Children() []Node {
	var res []Node
	res = append(res, n.Name)
	if n.Dependencies != nil {
		res = append(res, n.Dependencies)
	}

	if n.Input != nil {
		res = append(res, n.Input)
	}
	if n.Output != nil {
		res = append(res, n.Output)
	}

	if n.Block != nil {
		res = append(res, n.Block)
	}

	if n.ToDo != nil {
		res = append(res, n.ToDo)
	}

	if n.Definition != nil {
		res = append(res, n.Definition)
	}
	return res
}

type Stmts struct {
	node
	Statements []*Stmt `@@*`
}

func (n *Stmts) Children() []Node {
	if n == nil {
		return nil
	}

	var res []Node
	for _, statement := range n.Statements {
		res = append(res, statement)
	}

	return res
}

type Stmt struct {
	node

	IfStmt     *IfStmt     `@@`
	EachStmt   *EachStmt   `|@@`
	ToDo       *ToDo       `|@@`
	ReturnStmt *ReturnStmt `|@@`
	WhileStmt  *WhileStmt  `|@@`
	CallStmt   *CallStmt   `|@@`
	Block      *Stmts      `|"{" @@ "}"`
}

func (n *Stmt) Children() []Node {
	if n == nil {
		return nil
	}

	var res []Node
	if n.IfStmt != nil {
		res = append(res, n.IfStmt)
	}

	if n.EachStmt != nil {
		res = append(res, n.EachStmt)
	}

	if n.ToDo != nil {
		res = append(res, n.ToDo)
	}

	if n.ReturnStmt != nil {
		res = append(res, n.ReturnStmt)
	}

	if n.WhileStmt != nil {
		res = append(res, n.WhileStmt)
	}

	if n.CallStmt != nil {
		res = append(res, n.CallStmt)
	}

	if n.Block != nil {
		res = append(res, n.Block)
	}

	return res
}

type ReturnStmt struct {
	node
	Stmt *CallStmt `"gib" @@ "zurück"`
}

func (n *ReturnStmt) Children() []Node {
	return []Node{n.Stmt}
}

type WhileStmt struct {
	node
	Condition *CallStmt `"solange"  @@ `
	Body      *Stmt     `@@`
}

func (n *WhileStmt) Children() []Node {
	return []Node{n.Condition, n.Body}
}

type CallStmt struct {
	node
	Scribble string           `(@Text`
	Name     *TypeDeclaration `|@@)`
	Params   []*CallStmt      `( "(" (@@ ("," @@)*)? ")" )?`
}

func (n *CallStmt) Children() []Node {
	var res []Node
	if n.Name != nil {
		res = append(res, n.Name)
	}

	for _, param := range n.Params {
		res = append(res, param)
	}

	return res
}

type EachStmt struct {
	node

	Element  *TypeDeclaration `"für" ("jede"|"jedes"|"jeden") (@@)`
	Iterator *TypeDeclaration `"aus" @@`
	Body     *Stmt            `@@`
}

func (n *EachStmt) Children() []Node {
	return []Node{n.Element, n.Body}
}

type IfStmt struct {
	node

	Condition *CallStmt `"wenn" @@ "dann"`
	Body      *Stmt     `@@`
	Else      *Stmt     `("sonst" @@)?`
}

func (n *IfStmt) Children() []Node {
	var res []Node

	res = append(res, n.Condition)
	if n.Else != nil {
		res = append(res, n.Else)
	}

	return res
}

// Input defines a sum type of inputs, usually
// data types.
type Input struct {
	node
	Params []*TypeDeclaration `@@ ("und" @@)*`
}

func (n *Input) Children() []Node {
	if n == nil {
		return nil
	}

	var res []Node
	for _, param := range n.Params {
		res = append(res, param)
	}

	return nil
}

// Output defines a choice list of types.
type Output struct {
	node
	Params []*TypeDeclaration `@@ ("oder" @@)*`
}

func (n *Output) Children() []Node {
	if n == nil {
		return nil
	}

	var res []Node
	for _, param := range n.Params {
		res = append(res, param)
	}

	return nil
}
