package puml

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
)

func Func(r *resolver.Resolver, fun *parser.Function, flags RFlags) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	diag.DefaultTextAlignment = plantuml.DTACenter

	// stop recursion
	if done, ok := flags.Visited[fun]; ok && done {
		ac := &plantuml.ActLabelStmt{Name: fun.Name.Value}
		ac.Color = ColorTask
		ac.Name = bpmSym(bpmn_icon_subprocess_collapsed) + "\n"
		ac.Name += "//Rekursive Aufgabe//\n" + fun.Name.Value
		diag.Add(ac)
		return diag
	}

	flags.Visited[fun] = true

	// start
	mainType := flags.MainType == fun
	if mainType {
		diag.Add(&plantuml.ActStartStmt{})
	}

	// in
	paramsAsSplit := &plantuml.ActSplitStmt{}
	diag.Add(paramsAsSplit)
	for _, param := range fun.Params {
		paramsAsSplit.Stmts = append(paramsAsSplit.Stmts, newFunParam(r, param, false, mainType))
	}

	// body
	if fun.Body == nil || len(fun.Body.Children()) == 0 {
		// no body, draw as undefined primitive
		ac := &plantuml.ActLabelStmt{Name: fun.Name.Value}
		ac.Color = ColorTask
		ac.Name = bpmSym(bpmn_icon_task) + "\n"
		ac.Name += "//Allgemeine Aufgabe//\n" + fun.Name.Value
		diag.Add(ac)
		if mainType {
			extA, _ := parser.ParseExternalSystemAnnotation(fun.Parent().(*parser.TypeDefinition))
			if extA != nil {
				ac.Notes = append(ac.Notes, &plantuml.ActivityNote{
					Text: "Implementierung ist nicht\nBestandteil der Fachlichkeit.\nLeistung wird durch\nFremdsystem erbracht.",
				})
				ac.Color = ColorExternFunc
			}

		}
	} else {
		// has a body, thus need a partition structuring
		fn := &plantuml.ActPartitionStmt{Name: fun.Name.Value}
		diag.Add(fn)

		fn.Body = append(fn.Body, fromStmt(r, fun.Body, flags).ActivityStatements()...)
	}

	// out
	if fun.Return != nil {
		resultAsSplit := &plantuml.ActSplitStmt{}
		diag.Add(resultAsSplit)

		for _, param := range fun.Return.Params {
			stmts := plantuml.ActStmts{}
			//stmts = append(stmts, &plantuml.ActGotoLabel{Name: param.Name.String()}) TODO broken for white space and visual
			stmts = append(stmts, newFunParam(r, param, true, mainType))

			if mainType {
				stmts = append(stmts, &plantuml.ActKillStmt{})
			}

			resultAsSplit.Stmts = append(resultAsSplit.Stmts, stmts)

		}
	}

	return diag
}

func newFunParam(r *resolver.Resolver, dec *parser.TypeDeclaration, resultParam bool, resultAsFinalEvent bool) *plantuml.ActLabelStmt {
	eventName := typeDeclToLinkStr(r, dec)
	ac := &plantuml.ActLabelStmt{}
	ac.Color = ColorData
	if resultParam {
		evtA := getEventAnnotation(r, dec)
		if evtA != nil {
			if evtA.Out {
				ac.Name = bpmSym(bpmn_icon_end_event_message) + "\n"
				ac.Name += "//Domänenereignis ausgehend//\n" + eventName
				ac.Color = ColorEvent
			}
		}

		if ac.Name == "" {
			ac.Name = bpmSym(bpmn_icon_data_output) + "\n"
			ac.Name += "//ausgehende Daten//\n" + eventName
		}

		if resultAsFinalEvent {
			ac.Name = bpmSym(bpmn_icon_end_event_terminate) + "\n"
			if evtA != nil {
				ac.Name = bpmSym(bpmn_icon_end_event_message) + "\n"
				ac.Name += "//Domänenereignis ausgehend//\n" + eventName
				ac.Color = ColorEvent
			} else {
				ac.Name += "//Endergebnis//\n" + eventName
			}
		} else {
			ac.Color = ColorIntermediateFnResultEvent
		}

		if looksLikeError(r, dec) {
			ac.Color = ColorErrorEvent
			ac.Name = bpmSym(bpmn_icon_end_event_error) + "\n"
			ac.Name += "//Behandelter Fehler//\n" + eventName
		}
	} else {
		evtA := getEventAnnotation(r, dec)
		if evtA != nil {
			if evtA.In {
				ac.Name = bpmSym(bpmn_icon_receive) + "\n"
				ac.Name += "//Domänenereignis eingehend//\n" + eventName
				ac.Color = ColorEvent
			}
		}

		if ac.Name == "" {
			ac.Name = bpmSym(bpmn_icon_data_input) + "\n"
			ac.Name += "//eingehende Daten//\n" + eventName
		}

		typeDecl := r.Resolve(resolver.NewQualifiedNameFromLocalName(dec.Name))
		if len(typeDecl) > 0 {
			if _, ok := typeDecl[0].Type.(*parser.Function); ok {
				ac.Name = bpmSym(bpmn_icon_subprocess_collapsed) + "\n"
				if isExternalFunc(r, dec) {
					ac.Name += "//Fremdsystem//\n" + eventName
					ac.Color = ColorExternFunc
				} else {
					ac.Name += "//interne Abhängigkeit//\n" + eventName
					ac.Color = ColorInternFunc
				}

			}
		}

	}

	return ac
}

func isExternalFunc(r *resolver.Resolver, dec *parser.TypeDeclaration) bool {
	defs := r.ResolveLocalQualifier(dec.Name)
	if len(defs) == 0 {
		return false
	}

	if a, _ := parser.ParseExternalSystemAnnotation(defs[0]); a != nil {
		return true
	}

	return false
}

func getEventAnnotation(r *resolver.Resolver, dec *parser.TypeDeclaration) *parser.EventAnnotation {
	defs := r.ResolveLocalQualifier(dec.Name)
	if len(defs) == 0 {
		return nil
	}

	a, _ := parser.ParseEventAnnotation(defs[0])
	return a
}

func looksLikeError(r *resolver.Resolver, dec *parser.TypeDeclaration) bool {
	defs := r.ResolveLocalQualifier(dec.Name)
	if len(defs) == 0 {
		return false
	}

	if a, _ := parser.ParseErrorAnnotation(defs[0]); a != nil {
		return true
	}

	return false
}

func fromIfStmt(r *resolver.Resolver, ifStmt *parser.FnStmtIf, flags RFlags) *plantuml.ActIfStmt {
	stmt := &plantuml.ActIfStmt{
		Condition:    bpmSym(bpmn_icon_gateway_xor) + "\n" + ifStmt.Condition.Name.String() + "?", // TODO what about the params?
		PositiveText: "ja",
		NegativeText: "nein",
	}

	if ifStmt.Body != nil {
		stmt.PositiveStmt = append(stmt.PositiveStmt, fromStmt(r, ifStmt.Body, flags).ActivityStatements()...)
	}

	if ifStmt.Else != nil {
		stmt.NegativeStmt = append(stmt.NegativeStmt, fromStmt(r, ifStmt.Else, flags).ActivityStatements()...)
	}

	return stmt
}

func fromWhileStmt(r *resolver.Resolver, wStmt *parser.FnStmtWhile, flags RFlags) *plantuml.ActWhileStmt {
	stmt := &plantuml.ActWhileStmt{
		Condition:    bpmSym(bpmn_icon_loop_marker) + "\\n" + wStmt.Condition.Name.String() + "?", // this is different than IF !?
		PositiveText: "ja",
		NegativeText: "nein",
	}

	if wStmt.Body != nil {
		stmt.PositiveStmt = append(stmt.PositiveStmt, fromStmt(r, wStmt.Body, flags).ActivityStatements()...)
	}

	return stmt
}

func fromStmt(r *resolver.Resolver, stmt parser.FnStmt, flags RFlags) *plantuml.Diagram {
	switch t := stmt.(type) {
	case *parser.FnLitExpr:
		defs := r.ResolveLocalQualifier(t.Name)
		if len(defs) > 0 {
			if callFn, ok := defs[0].Type.(*parser.Function); ok {
				return Func(r, callFn, flags)

			}
		} else {
			diag := plantuml.NewDiagram()
			ac := &plantuml.ActLabelStmt{Name: t.Name.String()}
			ac.Color = ColorTaskUndefined
			ac.Name = bpmSym(bpmn_icon_task) + "\n"
			ac.Name += "//Undefinierte Aufgabe//\n" + t.Name.String()
			diag.Add(ac)
			return diag
		}

	case *parser.FnStmtBlock:
		diag := plantuml.NewDiagram()
		for _, statement := range t.Stmts.Statements {
			diag.Add(fromStmt(r, statement, flags).Renderables...)
		}

		return diag

	case *parser.FnStmtIf:
		return plantuml.NewDiagram().Add(fromIfStmt(r, t, flags))

	case *parser.FnStmtWhile:
		return plantuml.NewDiagram().Add(fromWhileStmt(r, t, flags))

	case *parser.FuncTypeRet:
		diag := plantuml.NewDiagram()
		resultAsSplit := &plantuml.ActSplitStmt{}
		diag.Add(resultAsSplit)
		for _, param := range t.Params {
			stmts := plantuml.ActStmts{}
			//resultAsSplit.Stmts = append(resultAsSplit.Stmts, &plantuml.ActGoto{Name: param.Name.String()}) //TODO: this looks entirely broken
			stmts = append(stmts, newFunParam(r, param, true, false))
			stmts = append(stmts, &plantuml.ActDetachStmt{})
			resultAsSplit.Stmts = append(resultAsSplit.Stmts, stmts)
		}
		return diag
	}

	return plantuml.NewDiagram().Add(&plantuml.ActLabelStmt{Name: fmt.Sprintf("unsupported type: %T", stmt)})
}
