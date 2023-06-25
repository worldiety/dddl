package puml

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"golang.org/x/exp/slog"
)

func Workflow(doc *parser.Doc, flow *parser.Workflow) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	ac := plantuml.NewActivity()
	diag.Add(ac)

	start := &plantuml.StartStmt{}
	ac.Stmts = append(ac.Stmts, &plantuml.Stmt{Start: start})
	if flow.Input != nil && len(flow.Input.Params) > 0 {
		note := &plantuml.ActivityNote{
			Color: "#aeebea",
			Text:  "Eingabe ist\n",
		}
		for i, param := range flow.Input.Params {
			note.Text += typeDeclToLinkStr(param) + "\n"
			if i < len(flow.Input.Params)-1 {
				note.Text += "und "
			}
		}

		start.Note = note
	}

	if flow.Dependencies != nil && len(flow.Dependencies.Params) > 0 {
		note := &plantuml.ActivityNote{
			Color:     "#fcd39a",
			Text:      "Hängt ab von\n",
			Direction: "left",
		}
		for i, param := range flow.Dependencies.Params {
			note.Text += typeDeclToLinkStr(param) + "\n"
			if i < len(flow.Dependencies.Params)-1 {
				note.Text += "und "
			}
		}

		ac.Stmts = append(ac.Stmts, &plantuml.Stmt{Note: note})
	}

	for _, statement := range flow.Block.Statements {
		pstate := fromStmt(statement)
		ac.Stmts = append(ac.Stmts, pstate)
	}

	// a common stop this does not make sense,
	// because we have to define the correct branches for each return

	fmt.Println(plantuml.String(diag))
	return diag
}

func fromStmt(stmt *parser.Stmt) *plantuml.Stmt {
	if stmt.ToDo != nil {
		return &plantuml.Stmt{Note: &plantuml.ActivityNote{Text: stmt.ToDo.Text.Text}}
	}
	if stmt.IfStmt != nil {
		return &plantuml.Stmt{IfStmt: fromIfStmt(stmt.IfStmt)}
	}

	if stmt.CallStmt != nil {
		return &plantuml.Stmt{State: fromCallStmt(stmt.CallStmt)}
	}

	if stmt.ReturnStmt != nil {
		return &plantuml.Stmt{Stop: fromReturnStmt(stmt.ReturnStmt)}
	}

	if stmt.WhileStmt != nil {
		return &plantuml.Stmt{While: fromWhileStmt(stmt.WhileStmt)}
	}

	if stmt.EachStmt != nil {
		return &plantuml.Stmt{While: fromEachStmt(stmt.EachStmt)}
	}

	if stmt.Block != nil {
		if len(stmt.Block.Statements) == 0 {
			return &plantuml.Stmt{State: plantuml.NewActivityState("Es wurde noch kein Zustand\noder Arbeitsschritt definiert").SetColor("#HotPink")}
		}

		if len(stmt.Block.Statements) == 1 {
			return fromStmt(stmt.Block.Statements[0])
		}

		partition := &plantuml.PartitionStmt{Body: &plantuml.Stmt{}}
		for _, statement := range stmt.Block.Statements {
			partition.Body.Block = append(partition.Body.Block, fromStmt(statement))
		}

		return &plantuml.Stmt{PartitionStmt: partition}
	}

	slog.Error("puml support: unsupported state", slog.Any("stmt", stmt))
	return &plantuml.Stmt{State: plantuml.NewActivityState("unsupported state")}
}

func fromReturnStmt(n *parser.ReturnStmt) *plantuml.StopStmt {
	stop := &plantuml.StopStmt{}
	if n.Stmt != nil {
		note := &plantuml.ActivityNote{
			Color: "#aeebb7",
			Text:  "Ausgabe ist ",
		}
		note.Text += param2String(n.Stmt)

		stop.Note = note
	}

	return stop
}

func fromCallStmt(n *parser.CallStmt) *plantuml.ActivityState {
	stateName := n.Scribble
	if n.Name != nil {
		stateName = TypeDeclToStr(n.Name)
	}

	ac := plantuml.NewActivityState(stateName)

	if len(n.Params) > 0 {
		note := &plantuml.ActivityNote{
			Color: "#aec8eb",
			Text:  "verwendet\n",
		}
		for i, param := range n.Params {
			note.Text += param2String(param) + "\n"
			if i < len(n.Params)-1 {
				note.Text += "und "
			}
		}

		ac.Note = note
	}

	return ac
}

func typeDeclToLinkStr(decl *parser.TypeDeclaration) string {
	tmp := decl.Name.Name
	if len(decl.Params) > 0 {
		tmp += "<"
		for i, param := range decl.Params {
			tmp += typeDeclToLinkStr(param)
			if i < len(decl.Params)-1 {
				tmp += ", "
			}
		}
		tmp += ">"
	} else {
		tmp = "[[#" + tmp + " " + tmp + "]]"
	}

	return tmp
}

func param2String(param *parser.CallStmt) string {
	tmp := param.Scribble
	if param.Name != nil {
		tmp = typeDeclToLinkStr(param.Name)
	}

	if len(param.Params) > 0 {
		tmp += "("
		for i, stmt := range param.Params {
			tmp += param2String(stmt)
			if i < len(param.Params)-1 {
				tmp += ", "
			}
		}
		tmp += ")"
	}

	return tmp
}

func fromIfStmt(ifStmt *parser.IfStmt) *plantuml.IfStmt {
	stmt := &plantuml.IfStmt{
		Condition:    fromCallStmt(ifStmt.Condition).Name + "?",
		PositiveText: "ja",
		NegativeText: "nein",
	}

	if ifStmt.Body != nil {
		stmt.PositiveStmt = fromStmt(ifStmt.Body)
	}

	if ifStmt.Else != nil {
		stmt.NegativeStmt = fromStmt(ifStmt.Else)
	}

	return stmt
}

func fromWhileStmt(n *parser.WhileStmt) *plantuml.WhileStmt {
	stmt := &plantuml.WhileStmt{
		Condition:    fromCallStmt(n.Condition).Name + "?",
		PositiveText: "ja",
		NegativeText: "nein",
	}

	if n.Body != nil {
		stmt.Body = fromStmt(n.Body)
	}

	return stmt
}

func fromEachStmt(n *parser.EachStmt) *plantuml.WhileStmt {
	elem := TypeDeclToStr(n.Element)
	it := TypeDeclToStr(n.Iterator)
	stmt := &plantuml.WhileStmt{
		Condition:    "aus " + it,
		PositiveText: elem,
		NegativeText: "alle Elemente verarbeitet",
	}

	if n.Body != nil {
		stmt.Body = fromStmt(n.Body)
	}

	return stmt
}
