package puml

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

func Workflow(doc *parser.Doc, flow *parser.Workflow) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	diag.DefaultTextAlignment = plantuml.DTACenter
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

	//log.Println(plantuml.String(diag))

	// pick first swimlane, if any
	_ = parser.Walk(flow, func(n parser.Node) error {
		if actor, ok := n.(*parser.ActorStmt); ok {
			stmt := &plantuml.Swimlane{Text: actor.ScribbleOrIdent.Value()}
			diag.Renderables = slices.Insert(diag.Renderables, 0, plantuml.Renderable(stmt))
			return fmt.Errorf("done")
		}

		return nil
	})

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
		return fromReturnStmt(stmt.ReturnStmt)
	}

	if stmt.ReturnErrorStmt != nil {
		return fromReturnErrorStmt(stmt.ReturnErrorStmt)
	}

	if stmt.WhileStmt != nil {
		return &plantuml.Stmt{While: fromWhileStmt(stmt.WhileStmt)}
	}

	if stmt.Event != nil {
		return &plantuml.Stmt{State: fromEventStmt(stmt.Event)}
	}

	if stmt.EventSent != nil {
		return &plantuml.Stmt{State: fromEventSentStmt(stmt.EventSent)}
	}

	if stmt.Actor != nil {
		pstmt := &plantuml.Stmt{Swimlane: fromActorStatement(stmt.Actor)}
		for _, statement := range stmt.Actor.Block.Statements {
			pstmt.Block = append(pstmt.Block, fromStmt(statement))
		}

		return pstmt
	}

	if stmt.Activity != nil {
		return &plantuml.Stmt{State: fromActivityStmt(stmt.Activity)}
	}

	if stmt.EachStmt != nil {
		return &plantuml.Stmt{While: fromEachStmt(stmt.EachStmt)}
	}

	if stmt.Block != nil {
		return &plantuml.Stmt{Block: fromStmts(stmt.Block)}
	}

	if stmt.Context != nil {
		return &plantuml.Stmt{PartitionStmt: fromContextStmt(stmt.Context)}
	}

	slog.Error("puml support: unsupported state", slog.Any("stmt", stmt))
	return &plantuml.Stmt{State: plantuml.NewActivityState("unsupported state")}
}

func fromStmts(n *parser.Stmts) []*plantuml.Stmt {
	var res []*plantuml.Stmt
	for _, statement := range n.Statements {
		res = append(res, fromStmt(statement))
	}

	return res
}

func fromContextStmt(n *parser.ContextStmt) *plantuml.PartitionStmt {
	partition := &plantuml.PartitionStmt{Body: &plantuml.Stmt{}}
	partition.Name = n.Name.Value
	if n.Block != nil {
		partition.Body = &plantuml.Stmt{Block: fromStmts(n.Block)}
	}

	return partition
}

func fromReturnStmt(n *parser.ReturnStmt) *plantuml.Stmt {
	eventName := ""
	if n.Stmt != nil {
		eventName = n.Stmt.Value
	}

	ac := plantuml.NewActivityState(eventName)
	ac.Color = "#ff992a"
	if eventName == "" {
		ac.Name = bpmSym(bpmn_icon_end_event_terminate) + "\n"
		ac.Name += "//Endereignis//\n" + eventName
	} else {
		ac.Name = bpmSym(bpmn_icon_end_event_message) + "\n"
		ac.Name += "//Nachrichten-Endereignis//\n" + eventName
	}

	return &plantuml.Stmt{Block: []*plantuml.Stmt{
		{State: ac},
		{Kill: &plantuml.KillStmt{}},
	}}

}

func fromReturnErrorStmt(n *parser.ReturnErrorStmt) *plantuml.Stmt {
	eventName := ""
	if n.Stmt != nil {
		eventName = n.Stmt.Value
	}

	ac := plantuml.NewActivityState(eventName)
	ac.Color = "#ec4d4e"
	if eventName == "" {
		ac.Name = bpmSym(bpmn_icon_end_event_cancel) + "\n"
		ac.Name += "//Abbruch//\n" + eventName
	} else {
		ac.Name = bpmSym(bpmn_icon_end_event_error) + "\n"
		ac.Name += "//Behandelter Fehler//\n" + eventName
	}

	return &plantuml.Stmt{Block: []*plantuml.Stmt{
		{State: ac},
		{Kill: &plantuml.KillStmt{}},
	}}

}

func bpmSym(symbol BpmnSymbol) string {
	return fmt.Sprintf("<size:25><font:bpmn><U+%s></font></size>", symbol)
}

func fromEventStmt(n *parser.EventStmt) *plantuml.ActivityState {
	eventName := n.Literal.Value
	ac := plantuml.NewActivityState(eventName)
	ac.Color = "#ff992a"
	ac.Name = bpmSym(bpmn_icon_receive) + "\n"
	ac.Name += "//Ereignis eingetreten//\n" + eventName
	return ac
}

func fromEventSentStmt(n *parser.EventSentStmt) *plantuml.ActivityState {
	eventName := n.Literal.Value
	ac := plantuml.NewActivityState(eventName)
	ac.Color = "#ff992a"
	ac.Name = bpmSym(bpmn_icon_end_event_message) + "\n"
	ac.Name += "//Ereignis versendet//\n" + eventName
	return ac
}

func fromActorStatement(n *parser.ActorStmt) *plantuml.Swimlane {
	actor := n.ScribbleOrIdent.Value()
	lane := &plantuml.Swimlane{Text: actor}
	return lane
}

func fromActivityStmt(n *parser.ActivityStmt) *plantuml.ActivityState {
	eventName := n.ScribbleOrIdent.Value()
	ac := plantuml.NewActivityState(eventName)
	ac.Color = "#3399fe"
	ac.Name = bpmSym(bpmn_icon_task) + "\n"
	ac.Name += "//Arbeitsschritt//\n" + eventName
	return ac
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
	tmp := decl.Name.Value
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
		Condition:    bpmSym(bpmn_icon_gateway_xor) + "\n" + ifStmt.Condition.Value + "?",
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
