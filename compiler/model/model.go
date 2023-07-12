// Package model contains a simplified and resolved representation of a DDD workspace.
// It is provided to simplify common template based code generators, which usually don't care about
// linting and resolving hassles.
package model

type Package struct {
	Comment     string
	Name        string
	ImportPath  string
	ChoiceTypes []*ChoiceType
	RecordTypes []*RecordType
	FuncTypes   []*FuncType
	Shared      bool
}

func (p *Package) TypeByName(name string) any {
	for _, funcType := range p.FuncTypes {
		if funcType.Name == name {
			return funcType
		}
	}

	for _, choiceType := range p.ChoiceTypes {
		if choiceType.Name == name {
			return choiceType
		}
	}

	for _, recordType := range p.RecordTypes {
		if recordType.Name == name {
			return recordType
		}
	}

	return nil
}

func (p *Package) HasType(name string) bool {
	return p.TypeByName(name) != nil
}

type ChoiceType struct {
	Parent  *Package
	Comment string
	Name    string
	Choices []*TypeDef
}

type RecordType struct {
	Parent  *Package
	Comment string
	Name    string
	Fields  []*Field
}

type Field struct {
	Parent  *RecordType
	Comment string
	Name    string
	Type    *TypeDef
}

type FuncType struct {
	Parent  *Package
	Comment string
	Name    string
	FuncDef *FuncDef
}

type QualifiedName struct {
	PackageImportPath string
	PackageName       string
	Name              string
}

func (q QualifiedName) IsUniverse() bool {
	return q.PackageName == "" && q.Name != ""
}

type TypeDef struct {
	Name      QualifiedName
	Parameter []*TypeDef
	FuncDef   *FuncDef // optionally we may be a functional definition
}

type FuncDef struct {
	Input  []*TypeDef
	Output *TypeDef // choice type
	Error  *TypeDef // choice type
}
