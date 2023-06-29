package plantuml

import "io"

type Activity struct {
	Stmts []*Stmt
}

func NewActivity() *Activity {
	return &Activity{}
}

func (a *Activity) Start() *Activity {
	a.Stmts = append(a.Stmts, &Stmt{Start: &StartStmt{}})
	return a
}

func (a *Activity) AddState(ac *ActivityState) *Activity {
	a.Stmts = append(a.Stmts, &Stmt{State: ac})
	return a
}

func (a *Activity) Render(wr io.Writer) error {
	for _, stmt := range a.Stmts {
		if err := stmt.Render(wr); err != nil {
			return err
		}
	}

	return nil
}

type Stmt struct {
	Start         *StartStmt
	While         *WhileStmt
	Stop          *StopStmt
	Kill          *KillStmt
	State         *ActivityState
	IfStmt        *IfStmt
	Block         []*Stmt
	PartitionStmt *PartitionStmt
	Note          *ActivityNote
	Swimlane      *Swimlane
}

func (n *Stmt) Render(wr io.Writer) error {

	if n.Start != nil {
		if err := n.Start.Render(wr); err != nil {
			return err
		}
	}

	if n.Stop != nil {
		if err := n.Stop.Render(wr); err != nil {
			return err
		}
	}

	if n.State != nil {
		if err := n.State.Render(wr); err != nil {
			return err
		}
	}

	if n.Swimlane != nil {
		if err := n.Swimlane.Render(wr); err != nil {
			return err
		}
	}

	if n.IfStmt != nil {
		if err := n.IfStmt.Render(wr); err != nil {
			return err
		}
	}

	if n.While != nil {
		if err := n.While.Render(wr); err != nil {
			return err
		}
	}

	if n.PartitionStmt != nil {
		if err := n.PartitionStmt.Render(wr); err != nil {
			return err
		}
	}

	for _, stmt := range n.Block {
		if err := stmt.Render(wr); err != nil {
			return err
		}
	}

	if n.Note != nil {
		if err := n.Note.Render(wr); err != nil {
			return err
		}
	}

	if n.Kill != nil {
		if err := n.Kill.Render(wr); err != nil {
			return err
		}
	}
	return nil
}

type ActivityNote struct {
	Direction string
	Color     string
	Text      string
}

func (n *ActivityNote) Render(wr io.Writer) error {
	dir := n.Direction
	if dir == "" {
		dir = "right"
	}
	w := strWriter{Writer: wr}
	w.Print("note ")
	w.Print(dir)
	if n.Color != "" {
		w.Print(" ")
		w.Print(n.Color)
	}
	w.Print("\n")
	w.Print(n.Text)
	w.Print("\n")
	w.Print("end note\n")

	return w.Err
}

type StartStmt struct {
	Note *ActivityNote
}

func (n *StartStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("start\n")
	if n.Note != nil {
		if err := n.Note.Render(wr); err != nil {
			return err
		}
	}

	return w.Err
}

type Swimlane struct {
	Text string
}

func (n *Swimlane) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("|%s|\n", n.Text)

	return w.Err
}

type StopStmt struct {
	Note *ActivityNote
}

func (n *StopStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("stop\n")
	if n.Note != nil {
		if err := n.Note.Render(wr); err != nil {
			return err
		}
	}

	return w.Err
}

type KillStmt struct {
}

func (n *KillStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("kill\n")

	return w.Err
}

type ActivityState struct {
	Color string
	Name  string
	Note  *ActivityNote
}

func NewActivityState(name string) *ActivityState {
	return &ActivityState{Name: name}
}

func (n *ActivityState) SetColor(c string) *ActivityState {
	n.Color = c
	return n
}

func (n *ActivityState) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	if n.Color != "" {
		w.Print(n.Color)
	}
	w.Printf(":%s;\n", n.Name)

	if n.Note != nil {
		if err := n.Note.Render(wr); err != nil {
			return err
		}
	}

	return w.Err
}

type WhileStmt struct {
	Condition    string
	PositiveText string
	NegativeText string
	Body         *Stmt
}

func (n *WhileStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("while (%s) ", n.Condition)
	if n.PositiveText != "" {
		w.Printf("is (%s)", n.PositiveText)
	}
	w.Print("\n")
	if n.Body != nil {
		if err := n.Body.Render(wr); err != nil {
			return err
		}
	}
	w.Print("endwhile")
	if n.PositiveText != "" {
		w.Printf(" (%s)", n.NegativeText)
	}
	w.Print("\n")
	return nil
}

type IfStmt struct {
	Condition    string
	PositiveText string
	PositiveStmt *Stmt
	NegativeText string
	NegativeStmt *Stmt
}

func NewIfStmt(condition string) *IfStmt {
	return &IfStmt{Condition: condition}
}

func (n *IfStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("if (%s) then (%s)\n", n.Condition, n.PositiveText)
	if n.PositiveStmt != nil {
		if err := n.PositiveStmt.Render(wr); err != nil {
			return err
		}
	}

	if n.NegativeStmt != nil {
		w.Printf("else (%s)\n", n.NegativeText)
		if err := n.NegativeStmt.Render(wr); err != nil {
			return err
		}
	}

	w.Print("endif\n")

	return w.Err
}

type PartitionStmt struct {
	Name string
	Body *Stmt
}

func (n *PartitionStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("partition %s {\n", n.Name)
	if n.Body != nil {
		if err := n.Body.Render(wr); err != nil {
			return err
		}
	}
	w.Print("}\n")

	return w.Err
}
