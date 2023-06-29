package parser

type Workflow struct {
	node
	KeywordWorkflow *KeywordWorkflow `@@`
	Name            *Ident           `@@ ("=" `
	ToDo            *ToDo            `@@? `
	Dependencies    *Input           `("Abhängigkeiten" ":" @@)?`
	Input           *Input           ` ("Eingabe" ":" @@`
	Output          *Output          `"Ausgabe" ":" @@)?`
	Block           *Stmts           `("Ablauf" "{" @@ "}")?`
	Definition      *Definition      `@@?)?`
}

func (n *Workflow) Children() []Node {
	return sliceOf(
		n.KeywordWorkflow,
		n.Name,
		n.ToDo,
		n.Dependencies,
		n.Input,
		n.Output,
		n.Block,
		n.Definition,
	)
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

type ContextStmt struct {
	node
	KeywordContext *KeywordContext `@@`
	Name           *Ident          `@@`
	Block          *Stmts          `"{" @@ "}"`
}

func (n *ContextStmt) Children() []Node {
	return sliceOf(n.KeywordContext, n.Name, n.Block)
}

type EventSentStmt struct {
	node
	KeywordEventSent *KeywordEventSent `@@`
	Literal          *Literal          `@@`
}

func (n *EventSentStmt) Children() []Node {
	return []Node{n.KeywordEventSent, n.Literal}
}

type EventStmt struct {
	node
	KeywordEvent *KeywordEvent `@@`
	Literal      *Literal      `@@`
}

func (n *EventStmt) Children() []Node {
	return []Node{n.KeywordEvent, n.Literal}
}

type ActorStmt struct {
	node
	KeywordEvent    *KeywordActor   `@@`
	ScribbleOrIdent *IdentOrLiteral `@@`
	Block           *Stmts          `"{" @@ "}"`
}

func (n *ActorStmt) Children() []Node {
	return []Node{n.KeywordEvent, n.ScribbleOrIdent, n.Block}
}

type ActivityStmt struct {
	node
	KeywordEvent    *KeywordActivity `@@`
	ScribbleOrIdent *IdentOrLiteral  `@@`
}

func (n *ActivityStmt) Children() []Node {
	return []Node{n.KeywordEvent, n.ScribbleOrIdent}
}

type Stmt struct {
	node

	IfStmt          *IfStmt          `@@`
	Event           *EventStmt       `|@@`
	EventSent       *EventSentStmt   `|@@`
	Activity        *ActivityStmt    `|@@`
	Actor           *ActorStmt       `|@@`
	Context         *ContextStmt     `|@@`
	EachStmt        *EachStmt        `|@@`
	ToDo            *ToDo            `|@@`
	ReturnStmt      *ReturnStmt      `|@@`
	ReturnErrorStmt *ReturnErrorStmt `|@@`
	WhileStmt       *WhileStmt       `|@@`
	CallStmt        *CallStmt        `|@@`
	Block           *Stmts           `|"{" @@ "}"`
}

func (n *Stmt) Children() []Node {
	return sliceOf(
		n.IfStmt,
		n.Event,
		n.EventSent,
		n.Activity,
		n.Actor,
		n.Context,
		n.EachStmt,
		n.ToDo,
		n.ReturnStmt,
		n.ReturnStmt,
		n.WhileStmt,
		n.CallStmt,
		n.Block,
	)
}

type ReturnStmt struct {
	node
	KeywordReturn *KeywordReturn `@@`
	Stmt          *Literal       `@@?`
}

type ReturnErrorStmt struct {
	node
	KeywordReturnError *KeywordReturnError `@@`
	Stmt               *Literal            `@@?`
}

func (n *ReturnStmt) Children() []Node {
	return sliceOf(n.KeywordReturn, n.Stmt)
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

	KeywordDecision *KeywordDecision `@@`
	KeywordIf       *KeywordIf       `@@`
	Condition       *Literal         `@@`
	KeywordThen     *KeywordThen     `@@`
	Body            *Stmt            `@@`
	KeywordElse     *KeywordElse     `( @@`
	Else            *Stmt            ` @@)?`
}

func (n *IfStmt) Children() []Node {
	return sliceOf(
		n.KeywordDecision,
		n.KeywordIf,
		n.Condition,
		n.KeywordThen,
		n.Body,
		n.KeywordElse,
		n.Else,
	)
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
