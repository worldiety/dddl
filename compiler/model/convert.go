package model

import (
	"fmt"
	"github.com/worldiety/dddl/parser"
	"golang.org/x/exp/slog"
	"strings"
)

// Convert uses some common patterns to create camel-case identifiers which are typical for Go or Java-like languages.
func Convert(ws *parser.Workspace) []*Package {
	var res []*Package
	freeElements := ws.CollectFreeDataOrWorkflow()
	if len(freeElements) > 0 {
		pkg := &Package{
			Comment: "Package shared contains context independent models and processes.",
			Name:    "shared",
			Shared:  true,
		}

		res = append(res, pkg)

		for _, node := range freeElements {
			if data, ok := node.(*parser.Data); ok {
				if data.IsChoiceType() {
					pkg.ChoiceTypes = append(pkg.ChoiceTypes, convertChoice(ws, pkg, data))
				} else {
					pkg.RecordTypes = append(pkg.RecordTypes, convertRecord(ws, pkg, data))
				}
			}

			if wf, ok := node.(*parser.Workflow); ok {
				isMetaProcess := false
				_ = parser.Walk(wf, func(n parser.Node) error {
					if _, ok := n.(*parser.ContextStmt); ok {
						isMetaProcess = true
						return fmt.Errorf("workflow represents meta cross-bounded context documentation")
					}
					return nil
				})
				if !isMetaProcess {
					pkg.FuncTypes = append(pkg.FuncTypes, convertWorkflow(ws, pkg, wf))
				}
			}
		}

	}

	for _, ctx := range ws.CollectContextChildren() {
		pkgName := strings.ToLower(makeIdentifier(ctx.Name))
		pkg := &Package{
			Comment: fmt.Sprintf("Package %s represents the bounded context '%s'.\n\n%s", pkgName, ctx.Name, getContextDoc(ws, ctx.Name)),
			Name:    pkgName,
		}
		res = append(res, pkg)

		for _, node := range ctx.Children {
			if data, ok := node.(*parser.Data); ok {
				if data.IsChoiceType() {
					pkg.ChoiceTypes = append(pkg.ChoiceTypes, convertChoice(ws, pkg, data))
				} else {
					pkg.RecordTypes = append(pkg.RecordTypes, convertRecord(ws, pkg, data))
				}
			}

			if wf, ok := node.(*parser.Workflow); ok {
				pkg.FuncTypes = append(pkg.FuncTypes, convertWorkflow(ws, pkg, wf))
			}
		}
	}

	return res
}

func convertWorkflow(ws *parser.Workspace, parent *Package, wf *parser.Workflow) *FuncType {
	typeName := makeUpIdentifier(wf.Name.Value)
	res := &FuncType{
		Parent:  parent,
		Comment: getDoc(wf),
		Name:    typeName,
		FuncDef: &FuncDef{},
	}

	// some language like Go require some kind of explicit thread context (time out, cancellation etc.)
	res.FuncDef.Input = append(res.FuncDef.Input, &TypeDef{Name: QualifiedName{Name: parser.UContext}})

	// grab the dependencies, which are always foreign workflows
	for _, activity := range wf.Dependencies() {
		decl := ws.ResolveTypeDeclaration(activity.ScribbleOrIdent.Name)
		referencedWorkflow, ok := decl.(*parser.Workflow)
		if !ok {
			slog.Error(fmt.Sprintf("workflow '%s' references sub-workflow '%s' which is not resolvable", wf.Name.Value, activity.ScribbleOrIdent.Value()))
			continue
		}

		// workflows are just defined function types
		res.FuncDef.Input = append(res.FuncDef.Input, &TypeDef{
			FuncDef: convertWorkflow(ws, parent, referencedWorkflow).FuncDef,
		})
	}

	// input or events are just defined function types, which return the stuff and a technical unspecified error
	for _, inputOrEvent := range wf.Inputs() {
		qualifier, ok := ws.Resolve(inputOrEvent.IdentOrLiteral().Name)
		if !ok {
			slog.Error(fmt.Sprintf("workflow '%s' references input or event '%s' which is not resolvable", wf.Name.Value, inputOrEvent.IdentOrLiteral().Name))
			continue
		}
		res.FuncDef.Input = append(res.FuncDef.Input, &TypeDef{
			FuncDef: &FuncDef{
				Input:  []*TypeDef{{Name: QualifiedName{Name: parser.UContext}}},
				Output: &TypeDef{Name: QualifiedName{PackageName: makePkgIdentifier(qualifier.Context.Name.Value), Name: makeUpIdentifier(qualifier.Name.Value)}},
				Error:  &TypeDef{Name: QualifiedName{Name: parser.UError}},
			},
		})
	}

	// output events are function calls which are by definition fire and forget into a message bus
	for _, event := range wf.OutputEvents() {
		qualifier, ok := ws.Resolve(event.Literal.Name)
		if !ok {
			slog.Error(fmt.Sprintf("workflow '%s' references an intermediate event '%s' which is not resolvable", wf.Name.Value, qualifier.Name.Value))
			continue
		}

		res.FuncDef.Input = append(res.FuncDef.Input, &TypeDef{
			FuncDef: &FuncDef{
				Input: []*TypeDef{
					{Name: QualifiedName{Name: parser.UContext}},
					{Name: QualifiedName{PackageName: makePkgIdentifier(qualifier.Context.Name.Value), Name: makeUpIdentifier(qualifier.Name.Value)}},
				},
				Error: &TypeDef{Name: QualifiedName{Name: parser.UError}},
			},
		})
	}

	// output are all return stmts which create an artificial choice type, just like sum of all error returns
	var outputDataTypes []*parser.Data
	for _, output := range wf.Output() {
		data := ws.ResolveData(output.Stmt.Name)
		if data == nil {
			slog.Error(fmt.Sprintf("workflow '%s' references an output data type '%s' which is not resolvable", wf.Name.Value, output.Stmt.Name))
			continue
		}

		outputDataTypes = append(outputDataTypes, data)
	}

	res.FuncDef.Output = declareNewChoiceType(parent, res.Name+"Result", outputDataTypes)

	// error are all exit points with errors which also create an artificial choice type, just like sum of all output types
	var errorDataType []*parser.Data
	for _, error := range wf.Errors() {
		data := ws.ResolveData(error.Stmt.Name)
		if data == nil {
			slog.Error(fmt.Sprintf("workflow '%s' references an error data type '%s' which is not resolvable", wf.Name.Value, error.Stmt.Name))
			continue
		}

		errorDataType = append(errorDataType, data)
	}

	res.FuncDef.Error = declareNewChoiceType(parent, res.Name+"Error", errorDataType)

	return res
}

func declareNewChoiceType(pkg *Package, name string, types []*parser.Data) *TypeDef {
	def := &TypeDef{Name: QualifiedName{PackageName: pkg.Name, Name: name}}
	if pkg.HasType(name) {
		return def
	}

	var choiceDefs []*TypeDef
	var choiceNames []string
	for _, data := range types {
		tname := makeUpIdentifier(data.Name.Value)
		choiceNames = append(choiceNames, tname)
		choiceDefs = append(choiceDefs, &TypeDef{
			Name: QualifiedName{PackageName: makePkgIdentifier(parser.ContextOf(data).Name.Value), Name: tname},
		})
	}

	pkg.ChoiceTypes = append(pkg.ChoiceTypes, &ChoiceType{
		Parent:  pkg,
		Comment: fmt.Sprintf("%s is the choice type of %s", name, strings.Join(choiceNames, "|")),
		Name:    name,
		Choices: choiceDefs,
	})

	return def
}

func convertChoice(ws *parser.Workspace, parent *Package, data *parser.Data) *ChoiceType {
	typeName := makeUpIdentifier(data.Name.Value)
	res := &ChoiceType{
		Parent:  parent,
		Comment: getDoc(data),
		Name:    typeName,
	}

	for _, tDef := range data.ChoiceTypes() {
		def := convertTypeDef(ws, tDef)
		res.Choices = append(res.Choices, def)
	}

	return res
}

func convertRecord(ws *parser.Workspace, parent *Package, data *parser.Data) *RecordType {
	typeName := makeUpIdentifier(data.Name.Value)
	res := &RecordType{
		Parent:  parent,
		Comment: getDoc(data),
		Name:    typeName,
	}

	for _, tDef := range data.FieldTypes() {
		def := convertTypeDef(ws, tDef)
		res.Fields = append(res.Fields, &Field{
			Parent: res,
			Name:   makeUpIdentifier(tDef.Name.Value),
			Type:   def,
		})
	}

	return res
}

func convertTypeDef(ws *parser.Workspace, tDef *parser.TypeDef) *TypeDef {
	def := &TypeDef{}
	if tDef.Name.IsUniverse() {
		def.Name = QualifiedName{Name: tDef.Name.NormalizeUniverse()}

	} else {
		qualifier, ok := ws.Resolve(tDef.Name)
		if !ok {
			def.Name = QualifiedName{Name: parser.UAny}

		} else {
			pkgName := strings.ToLower(makeIdentifier(qualifier.Context.Name.Value))
			def.Name = QualifiedName{
				PackageName: pkgName,
				Name:        qualifier.Name.Value,
			}
		}
	}

	for _, param := range tDef.Params {
		def.Parameter = append(def.Parameter, convertTypeDef(ws, param))
	}

	return def
}

// getContextDoc joins all available definitions and todos into a big ball of text.
func getContextDoc(ws *parser.Workspace, name string) string {
	var sb strings.Builder
	for _, context := range ws.Contexts() {
		if context.Name.Value == name {
			sb.WriteString(getDoc(context))
		}
	}

	return sb.String()
}

func getDoc(def parser.Defineable) string {
	var sb strings.Builder
	if defText := parser.TextOf(def.GetDefinition()); defText != "" {
		sb.WriteString(defText)
	}

	if todoText := parser.TextOf(def.GetToDo()); todoText != "" {
		if sb.Len() > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString("# TODO\n\n")
		sb.WriteString(todoText)
	}

	return sb.String()
}
