package plantuml

import (
	"io"
	"regexp"
	"strings"
)

type States struct {
	states      map[string]*State
	transitions []*StateTransition
}

func NewStates() *States {
	return &States{states: map[string]*State{}}
}

func (n *States) Transition(t *StateTransition) {
	if t.From.Id == "" {
		t.From.Id = MakeId(t.From.Title)
	}

	if t.To.Id == "" {
		t.To.Id = MakeId(t.To.Title)
	}

	if _, ok := n.states[t.From.Id]; !ok {
		n.states[t.From.Id] = t.From
	}

	if _, ok := n.states[t.To.Id]; !ok {
		n.states[t.To.Id] = t.To
	}

	n.transitions = append(n.transitions, t)
}

func (n *States) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	for _, state := range n.states {
		if state.Id == "[*]" {
			continue
		}

		if err := state.Render(wr); err != nil {
			return err
		}
	}

	for _, transition := range n.transitions {
		if err := transition.Render(wr); err != nil {
			return err
		}
	}

	return w.Err
}

type StateStereotype string

const (
	StateJoin      StateStereotype = "join"
	StateFork      StateStereotype = "fork"
	StateInputPin  StateStereotype = "inputPin"
	StateOutputPin StateStereotype = "outputPin"
	StateStart     StateStereotype = "start"
	StateEnd       StateStereotype = "end"
	StateChoice    StateStereotype = "choice"
)

type State struct {
	Id          string
	Title       string
	Description string
	Color       string
	Stereotype  StateStereotype
	States      *States
}

func (n *State) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf(`state "%s" as %s`, escapeP(n.Title), n.Id)
	if n.Stereotype != "" {
		w.Printf(" <<%s>> ", n.Stereotype)
	}

	if n.Color != "" {
		w.Printf(" %s ", n.Color)
	}

	if n.States != nil {
		w.Printf(" {\n")
		if err := n.Render(wr); err != nil {
			return err
		}
		w.Printf("}\n")
	}

	w.Print("\n")
	if n.Description != "" {
		w.Printf("%s : %s\n", n.Id, n.Description)
	}

	return w.Err
}

func NewStartState() *State {
	return &State{Id: "[*]"}
}

type StateTransition struct {
	From *State
	To   *State
	Text string
}

func (n *StateTransition) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("%s --> %s", n.From.Id, n.To.Id)
	if n.Text != "" {
		w.Printf(" : %s", n.Text)
	}
	w.Printf("\n")
	return w.Err
}

var regexNonWord = regexp.MustCompile(`[^\w]+`)

func MakeId(s string) string {
	s = strings.ToLower(s)
	return regexNonWord.ReplaceAllString(s, "_")
}
