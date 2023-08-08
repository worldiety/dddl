package parser

type UniverseName string

func (n UniverseName) IsList() bool {
	return n == UList || n == UListDE
}

func (n UniverseName) IsSet() bool {
	return n == USet || n == USetDE
}

func (n UniverseName) IsMap() bool {
	return n == UMap || n == UMapDE
}

func (n UniverseName) IsNumber() bool {
	return n == UNumber || n == UNumberDE
}

func (n UniverseName) IsFloat() bool {
	return n == UFloat || n == UFloatDE
}

func (n UniverseName) IsInt() bool {
	return n == UInt || n == UIntDE
}

func (n UniverseName) IsAny() bool {
	return n == UAny
}

func (n UniverseName) IsFunc() bool {
	return n == UFunc
}

func (n UniverseName) IsString() bool {
	return n == UString || n == UStringDE
}

func (n UniverseName) NormalizeUniverse() string {
	if n.IsList() {
		return UList
	}

	if n.IsSet() {
		return USet
	}

	if n.IsMap() {
		return UMap
	}

	if n.IsNumber() {
		return UNumber
	}

	if n.IsFloat() {
		return UFloat
	}

	if n.IsInt() {
		return UInt
	}

	if n.IsAny() {
		return UAny
	}

	if n.IsFunc() {
		return UFunc
	}

	if n.IsString() {
		return UString
	}

	return string(n)
}

const (
	UList     = "List"
	UListDE   = "Liste"
	USet      = "Set"
	USetDE    = "Menge"
	UMap      = "Map"
	UMapDE    = "Zuordnung"
	UString   = "String"
	UStringDE = "Text"
	UNumber   = "Number"
	UNumberDE = "Zahl"
	UInt      = "Integer"
	UIntDE    = "Ganzzahl"
	UFloat    = "Float"
	UFloatDE  = "Gleitkommazahl"

	// TODO: good or bad idea?
	UAny     = "any"
	UFunc    = "func"
	UError   = "error"
	UContext = "context"
)

func (n UniverseName) IsUniverse() bool {
	// Boolean is not defined, we encourage to use a choicetype
	switch n {
	//list, set, map, string, int, float
	case UListDE, USetDE, UMapDE, UStringDE, UNumberDE, UIntDE, UFloatDE:
		return true
	case UList, USet, UMap, UString, UNumber, UInt, UFloat:
		return true
	case UAny, UFunc, UError, UContext:
		return true
	default:
		return false
	}
}
