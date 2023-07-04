package parser

type Workflow struct {
	node
	KeywordWorkflow *KeywordWorkflow `@@`
	Name            *Ident           `@@ ( "{" `
	ToDo            *ToDo            `@@? `
	Definition      *Definition      `@@?`

	Block *Stmts `@@?  "}" )?`
}

func WorkflowOf(root Node) *Workflow {
	for root != nil {
		if wf, ok := root.(*Workflow); ok {
			return wf
		}
		root = root.Parent()
	}

	return nil
}

func (n *Workflow) DeclaredName() *Ident {
	return n.Name
}

func (n *Workflow) Qualifier() Qualifier {
	return Qualifier{
		Context: n.Parent().(*Context),
		Name:    n.Name,
	}
}

func (n *Workflow) Children() []Node {
	return sliceOf(
		n.KeywordWorkflow,
		n.Name,
		n.ToDo,
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
	Literal          *IdentOrLiteral   `@@`
}

func (n *EventSentStmt) Children() []Node {
	return sliceOf(n.KeywordEventSent, n.Literal)
}

type EventStmt struct {
	node
	KeywordEvent *KeywordEvent   `@@`
	Literal      *IdentOrLiteral `@@`
}

func (n *EventStmt) Children() []Node {
	return sliceOf(n.KeywordEvent, n.Literal)
}

type ActorStmt struct {
	node
	KeywordEvent    *KeywordActor   `@@`
	ScribbleOrIdent *IdentOrLiteral `@@`
	Block           *Stmts          `"{" @@ "}"`
}

func (n *ActorStmt) Children() []Node {
	return sliceOf(n.KeywordEvent, n.ScribbleOrIdent, n.Block)
}

type ViewStmt struct {
	node
	KeywordView     *KeywordView    `@@`
	ScribbleOrIdent *IdentOrLiteral `@@`
}

func (n *ViewStmt) Children() []Node {
	return sliceOf(n.KeywordView, n.ScribbleOrIdent)
}

type OutputStmt struct {
	node
	KeywordOutput   *KeywordOutput  `@@`
	ScribbleOrIdent *IdentOrLiteral `@@`
}

func (n *OutputStmt) Children() []Node {
	return sliceOf(n.KeywordOutput, n.ScribbleOrIdent)
}

type InputStmt struct {
	node
	KeywordInput    *KeywordInput   `@@`
	ScribbleOrIdent *IdentOrLiteral `@@`
}

func (n *InputStmt) Children() []Node {
	return sliceOf(n.KeywordInput, n.ScribbleOrIdent)
}

type ActivityStmt struct {
	node
	KeywordEvent    *KeywordActivity `@@`
	ScribbleOrIdent *IdentOrLiteral  `@@ ( "{"`
	ViewStmt        *ViewStmt        `@@?`
	InputStmt       *InputStmt       `@@?`
	OutputStmt      *OutputStmt      `@@? "}")?`
}

func (n *ActivityStmt) Children() []Node {
	return sliceOf(n.KeywordEvent, n.ScribbleOrIdent, n.ViewStmt, n.InputStmt, n.OutputStmt)
}

type Stmt struct {
	node

	IfStmt          *IfStmt          `@@`
	Event           *EventStmt       `|@@`
	EventSent       *EventSentStmt   `|@@`
	Activity        *ActivityStmt    `|@@`
	Actor           *ActorStmt       `|@@`
	Context         *ContextStmt     `|@@`
	ToDo            *ToDo            `|@@`
	ReturnStmt      *ReturnStmt      `|@@`
	ReturnErrorStmt *ReturnErrorStmt `|@@`
	WhileStmt       *WhileStmt       `|@@`
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
		n.ToDo,
		n.ReturnStmt,
		n.ReturnErrorStmt,
		n.WhileStmt,
		n.Block,
	)
}

type ReturnStmt struct {
	node
	KeywordReturn *KeywordReturn  `@@`
	Stmt          *IdentOrLiteral `@@?`
}

func (n *ReturnStmt) Children() []Node {
	return sliceOf(n.KeywordReturn, n.Stmt)
}

type ReturnErrorStmt struct {
	node
	KeywordReturnError *KeywordReturnError `@@`
	Stmt               *IdentOrLiteral     `@@?`
}

func (n *ReturnErrorStmt) Children() []Node {
	return sliceOf(n.KeywordReturnError, n.Stmt)
}

type WhileStmt struct {
	node
	KeywordWhile *KeywordWhile `@@`
	Condition    *Literal      ` @@ `
	Body         *Stmt         `@@`
}

func (n *WhileStmt) Children() []Node {
	return sliceOf(n.KeywordWhile, n.Condition, n.Body)
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
	Params []*TypeDef `@@ ("und" @@)*`
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
	Params []*TypeDef `@@ ("oder" @@)*`
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
