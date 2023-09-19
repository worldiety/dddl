package parser

import "fmt"

func parseAnnotations(definition *TypeDefinition) []TypedAnnotation {
	panic("todo")
}

type EventAnnotation struct {
	TypedAnnotation
	In  bool
	Out bool
}

// ParseEventAnnotation supports variants like:
// * @Ereignis
// * @Ereignis(eingehend,ausgehend)
// * @event(incoming,outgoing)
// Attention: even though error is nil, EventAnnotation may still be nil, because the absence is valid.
func ParseEventAnnotation(definition *TypeDefinition) (*EventAnnotation, error) {
	annotation, err := definition.ExpectOneOrNoneOf("Ereignis", "event")
	if err != nil {
		return nil, err
	}

	if annotation == nil {
		return nil, nil
	}

	switch definition.Type.(type) {
	case *Struct, *Choice, *Type, *Alias:
		// allowed
	default:
		return nil, fmt.Errorf("annotation %s is not allowed on type %s", annotation.Name.Value, definition.Type.GetName())
	}

	if err := annotation.ExpectKeysOf("eingehend", "ausgehend", "incoming", "outgoing"); err != nil {
		return nil, err
	}

	inVal, inOk := annotation.FirstValue("eingehend", "incoming")
	if inVal != "" {
		return nil, fmt.Errorf("%s.%s must not provide a value (%s)", annotation.Name.Value, "incoming", inVal)
	}

	outVal, outOk := annotation.FirstValue("ausgehend", "outgoing")
	if outVal != "" {
		return nil, fmt.Errorf("%s.%s must not provide a value (%s)", annotation.Name.Value, "outgoing", inVal)
	}

	return &EventAnnotation{
		In:  inOk,
		Out: outOk,
	}, nil
}

type ErrorAnnotation struct {
}

// ParseErrorAnnotation supports variants like:
// * @error
// * @Fehler
// Attention: even though error is nil, ErrorAnnotation may still be nil, because the absence is valid.
func ParseErrorAnnotation(definition *TypeDefinition) (*ErrorAnnotation, error) {
	annotation, err := definition.ExpectOneOrNoneOf("Fehler", "error")
	if err != nil {
		return nil, err
	}

	if annotation == nil {
		return nil, nil
	}

	switch definition.Type.(type) {
	case *Struct, *Choice, *Type, *Alias:
		// allowed
	default:
		return nil, fmt.Errorf("annotation %s is not allowed on type %s", annotation.Name.Value, definition.Type.GetName().Value)
	}

	if err := annotation.ExpectEmpty(); err != nil {
		return nil, err
	}

	return &ErrorAnnotation{}, nil
}

type ExternalSystemAnnotation struct {
}

// ParseExternalSystemAnnotation supports variants like:
// * @external
// * @Fremdsystem
// Attention: even though error is nil, ExternalSystemAnnotation may still be nil, because the absence is valid.
func ParseExternalSystemAnnotation(definition *TypeDefinition) (*ExternalSystemAnnotation, error) {
	annotation, err := definition.ExpectOneOrNoneOf("Fremdsystem", "external")
	if err != nil {
		return nil, err
	}

	if annotation == nil {
		return nil, nil
	}

	switch t := definition.Type.(type) {
	case *Function:
		if t.Body != nil {
			return nil, fmt.Errorf("external function '%s' must not have a body", t.Name.Value)
		}

		// allowed
	default:
		return nil, fmt.Errorf("annotation %s is not allowed on type %s", annotation.Name.Value, definition.Type.GetName().Value)
	}

	if err := annotation.ExpectEmpty(); err != nil {
		return nil, err
	}

	return &ExternalSystemAnnotation{}, nil
}
