package parser

import (
	"fmt"
	"golang.org/x/exp/slices"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

// parseAnnotations always returns a non-nil slice.
func parseAnnotations(tDef *TypeDefinition) []TypedAnnotation {
	tmp := make([]TypedAnnotation, 0, len(tDef.Annotations)) // optimize and allocate on 0 size
	if a := parseEventAnnotation(tDef); a != nil {
		tmp = append(tmp, a)
	}

	if a := parseErrorAnnotation(tDef); a != nil {
		tmp = append(tmp, a)
	}

	if a := parseExternalSystemAnnotation(tDef); a != nil {
		tmp = append(tmp, a)
	}

	if a := parseActor(tDef); a != nil {
		tmp = append(tmp, a)
	}

	if a := parseRoleAnnotation(tDef); a != nil {
		tmp = append(tmp, a)
	}

	if a := parseWorkPackageAnnotation(tDef); a != nil {
		tmp = append(tmp, a)
	}

	return tmp
}

type EventAnnotation struct {
	TypedAnnotation
	In  bool
	Out bool
}

// parseEventAnnotation supports variants like:
// * @Ereignis
// * @Ereignis(eingehend,ausgehend)
// * @event(incoming,outgoing)
func parseEventAnnotation(tDef *TypeDefinition) *EventAnnotation {
	switch tDef.Type.(type) {
	case *Struct, *Choice, *Type, *Alias:
		// allowed
	default:
		return nil
	}

	for _, annotation := range tDef.Annotations {
		if v := annotation.Name.Value; v == "Ereignis" || v == "event" {
			_, inOk := annotation.FirstValue("eingehend", "incoming")
			_, outOk := annotation.FirstValue("ausgehend", "outgoing")

			return &EventAnnotation{
				In:  inOk,
				Out: outOk,
			}

		}
	}

	return nil
}

type ErrorAnnotation struct {
	TypedAnnotation
}

// ParseErrorAnnotation supports variants like:
// * @error
// * @Fehler
func parseErrorAnnotation(tDef *TypeDefinition) *ErrorAnnotation {
	switch tDef.Type.(type) {
	case *Struct, *Choice, *Type, *Alias:
		// allowed
	default:
		return nil
	}

	for _, annotation := range tDef.Annotations {
		if v := annotation.Name.Value; v == "Fehler" || v == "error" {

			return &ErrorAnnotation{}

		}
	}

	return nil
}

type ExternalSystemAnnotation struct {
	TypedAnnotation
}

// parseExternalSystemAnnotation supports variants like:
// * @external
// * @Fremdsystem
func parseExternalSystemAnnotation(tDef *TypeDefinition) *ExternalSystemAnnotation {
	switch tDef.Type.(type) {
	case *Function:
		// allowed
	default:
		return nil
	}

	for _, annotation := range tDef.Annotations {
		if v := annotation.Name.Value; v == "Fremdsystem" || v == "external" {
			return &ExternalSystemAnnotation{}
		}
	}

	return nil
}

type TaskActors struct {
	TypedAnnotation
	Roles []string
}

// parseActor supports variants like:
// * @Rollen
// * @roles
// * @Rollen("a")
// * @Rollen("a","b","c")
func parseActor(tDef *TypeDefinition) *TaskActors {
	switch tDef.Type.(type) {
	case *Function:
		// allowed
	default:
		return nil
	}

	for _, annotation := range tDef.Annotations {
		if v := annotation.Name.Value; v == "Akteur" || v == "actor" {
			roles := &TaskActors{}
			for _, kv := range annotation.KeyValues {
				roles.Roles = append(roles.Roles, kv.Key.Value)
			}

			slices.Sort(roles.Roles)
			roles.Roles = slices.Compact(roles.Roles)

			return roles
		}
	}

	return nil
}

type RoleAnnotation struct {
	TypedAnnotation
}

// ParseErrorAnnotation supports variants like:
// * @role
// * @Rolle
func parseRoleAnnotation(tDef *TypeDefinition) *RoleAnnotation {
	switch tDef.Type.(type) {
	case *Struct, *Choice, *Type, *Alias:
		// allowed
	default:
		return nil
	}

	for _, annotation := range tDef.Annotations {
		if v := annotation.Name.Value; v == "Rolle" || v == "role" {
			return &RoleAnnotation{}
		}
	}

	return nil
}

type WorkPackageAnnotation struct {
	TypedAnnotation
	Name     string
	Duration time.Duration
	Requires []string
}

func (a *WorkPackageAnnotation) GetName() string {
	if a == nil {
		return ""
	}

	return a.Name
}

func (a *WorkPackageAnnotation) GetDuration() string {
	if a == nil {
		return ""
	}

	if a.Duration == 0 {
		return ""
	}

	v := a.Duration.Hours() / 8
	if v < 1 {
		return "< 1 PT"
	}

	v = math.Ceil(v)

	return fmt.Sprintf("~%d PT", int(v))
}

func (a *WorkPackageAnnotation) GetRequires() []string {
	if a == nil {
		return nil
	}

	return a.Requires
}

// parseActor supports variants like:
// * @Arbeitspaket(Name="Version 1.1", Aufwand="0.5d", benötigt="Version 1.0")
// * @workpackage(name="Milestone 5", duration="30m", depends="Version 1.0")
func parseWorkPackageAnnotation(tDef *TypeDefinition) *WorkPackageAnnotation {
	for _, annotation := range tDef.Annotations {
		if v := annotation.Name.Value; v == "Arbeitspaket" || v == "workpackage" {
			pkg := &WorkPackageAnnotation{}
			if duration, ok := annotation.FirstValue("Aufwand", "duration"); ok {
				parsedDuration, err := parseXDuration(duration)
				if err != nil {
					log.Printf("cannot parse duration from annotation (%s): %v", duration, err)
				}

				pkg.Duration = parsedDuration
			}

			for _, kv := range annotation.KeyValues {
				if kv.Key.Value == "benötigt" || kv.Key.Value == "depends" {
					pkg.Requires = append(pkg.Requires, kv.Value.Value)
				}
			}

			name, _ := annotation.FirstValue("Name", "name")
			pkg.Name = name

			slices.Sort(pkg.Requires)
			pkg.Requires = slices.Compact(pkg.Requires)

			return pkg
		}
	}

	return nil
}

func parseXDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		vFloat, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}

		return time.Duration(vFloat * 8 * float64(time.Hour)), nil
	}

	if strings.HasSuffix(s, "y") {
		vFloat, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}

		return time.Duration(vFloat * 8 * 365 * float64(time.Hour)), nil
	}

	return time.ParseDuration(s)
}
