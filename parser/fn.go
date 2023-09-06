package parser

import "github.com/alecthomas/participle/v2/lexer"

type KeywordFn struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("task" | "Aufgabe")`
}

func (n *KeywordFn) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type KeywordRet struct {
	node
	Tokens  []lexer.Token
	Keyword string `@("->")`
}

func (n *KeywordRet) EndPosition() lexer.Position {
	return n.relocateEndPos(n.Tokens)
}

type FuncTypeRet struct {
	node
	KeywordRet *KeywordRet        `@@`
	Params     []*TypeDeclaration `(@@ | "(" @@ ("," @@)* ")" )`
}

func (n *FuncTypeRet) Children() []Node {
	res := sliceOf(n.KeywordRet)
	for _, param := range n.Params {
		res = append(res, param)
	}

	return res
}

func (*FuncTypeRet) fnStmt() {}

type Function struct {
	node
	KeywordFn *KeywordFn `@@`
	Name      *Name      `@@`

	Params []*TypeDeclaration `(@@ | "(" @@ ("," @@)* ")" )?`
	// Return may be nil, for void return
	Return *FuncTypeRet `@@?`
	// Body may be nil for stubs or unknown branch-behavior
	Body *FnStmtBlock `( @@ )?`
}

func (n *Function) IsExternalSystem() bool {
	a, _ := ParseExternalSystemAnnotation(n.Parent().(*TypeDefinition))
	return a != nil
}

func (n *Function) GetKeyword() string {
	return n.KeywordFn.Keyword
}

func (n *Function) Children() []Node {
	res := sliceOf(n.KeywordFn, n.Name, n.Return, n.Body)
	for _, param := range n.Params {
		res = append(res, param)
	}
	return res
}

func (n *Function) GetName() *Name {
	return n.Name
}

func (*Function) namedType() {}

type TypeDeclaration struct {
	node
	Name *QualifiedName `@@`
	// Type parameters, like "Map[Key, Value]"
	Params []*TypeDeclaration `("[" @@ ("," @@)* "]" )?`
	// Choice is nil if "|TypeDecl" is following as a union, like "String | none". This allows anonymous choice types.
	Choice   *TypeDeclaration `("|" @@)?`
	Optional bool             `@"?"?`
}

func (n *TypeDeclaration) Children() []Node {
	res := sliceOf(n.Name, n.Choice)
	for _, param := range n.Params {
		res = append(res, param)
	}

	return res
}

type FnStmts struct {
	node
	Statements []FnStmt `@@*`
}

func (n *FnStmts) Children() []Node {
	var res []Node
	for _, statement := range n.Statements {
		res = append(res, statement)
	}

	return res
}

type FnStmt interface {
	Node
	fnStmt()
}

type FnStmtIf struct {
	node
	KeywordIf *KeywordIf `@@`
	Condition *FnLitExpr `@@`
	Body      FnStmt     `@@`
	// may be nil
	KeywordElse *KeywordElse `( @@`
	// may be nil
	Else FnStmt ` @@)?`
}

func (n *FnStmtIf) Children() []Node {
	return sliceOf(n.KeywordIf, n.Condition, n.Body, n.KeywordElse, n.Else)
}

func (*FnStmtIf) fnStmt() {}

type FnStmtWhile struct {
	node
	KeywordWhile *KeywordWhile `@@`
	Condition    *FnLitExpr    `@@`
	Body         FnStmt        `@@`
}

func (n *FnStmtWhile) Children() []Node {
	return sliceOf(n.KeywordWhile, n.Condition, n.Body)
}

func (*FnStmtWhile) fnStmt() {}

type FnStmtBlock struct {
	node
	Stmts *FnStmts `"{" @@ "}"`
}

func (n *FnStmtBlock) Children() []Node {
	return sliceOf(n.Stmts)
}

func (*FnStmtBlock) fnStmt() {}

type FnLitExpr struct {
	node
	Name   *QualifiedName     `@@`
	Params []*TypeDeclaration `("(" (@@ ("," @@)*)? ")" )?`
}

func (n *FnLitExpr) Children() []Node {
	res := sliceOf(n.Name)
	for _, param := range n.Params {
		res = append(res, param)
	}

	return res
}

func (*FnLitExpr) fnStmt() {}
