package plantuml

import "io"

type ActivityStatement interface {
	acStmt() // marker method
	Renderable
}

type ActLabelStmt struct {
	Color string
	Name  string
	Notes []*ActivityNote
}

func (n *ActLabelStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	if n.Color != "" {
		w.Print(n.Color)
	}
	w.Printf(":%s;\n", n.Name)

	for _, note := range n.Notes {
		_ = note.Render(w)
	}

	return w.Err
}

func (n *ActLabelStmt) acStmt() {}

type ActSplitStmt struct {
	Stmts []ActivityStatement
}

func (n *ActSplitStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	if len(n.Stmts) < 1 {
		return nil
	}

	if len(n.Stmts) == 1 {
		_ = n.Stmts[0].Render(w)
		return w.Err
	}

	w.Printf("split\n")
	for i, stmt := range n.Stmts {
		_ = stmt.Render(w)

		if i < len(n.Stmts)-1 {
			w.Printf("split again\n")
		}
	}
	w.Printf("end split\n")

	return w.Err
}

func (n *ActSplitStmt) acStmt() {}

type ActStmts []ActivityStatement

func (n ActStmts) acStmt() {}

func (n ActStmts) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	for _, statement := range n {
		_ = statement.Render(w)
		w.Print("\n")
	}

	return w.Err
}

type ActPartitionStmt struct {
	Name string
	Body ActStmts
}

func (n *ActPartitionStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("partition %s {\n", n.Name)
	if n.Body != nil {
		_ = n.Body.Render(wr)
	}
	w.Print("}\n")

	return w.Err
}

func (n *ActPartitionStmt) acStmt() {}

type ActStartStmt struct {
	Note *ActivityNote
}

func (n *ActStartStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("start\n")
	if n.Note != nil {
		if err := n.Note.Render(wr); err != nil {
			return err
		}
	}

	return w.Err
}

func (n *ActStartStmt) acStmt() {}

type ActKillStmt struct {
}

func (n *ActKillStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("kill\n")

	return w.Err
}

func (n *ActKillStmt) acStmt() {}

type ActDetachStmt struct {
}

func (n *ActDetachStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("detach\n")

	return w.Err
}

func (n *ActDetachStmt) acStmt() {}

type ActWhileStmt struct {
	Condition    string
	PositiveText string
	PositiveStmt ActStmts
	NegativeText string
	Body         *Stmt
}

func (n *ActWhileStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("while (%s) ", n.Condition)
	if n.PositiveText != "" {
		w.Printf("is (%s)", n.PositiveText)
	}
	w.Print("\n")

	if len(n.PositiveStmt) != 0 {
		_ = n.PositiveStmt.Render(w)
	}

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

func (n *ActWhileStmt) acStmt() {}

type ActIfStmt struct {
	Condition    string
	PositiveText string
	PositiveStmt ActStmts
	NegativeText string
	NegativeStmt ActStmts
}

func (n *ActIfStmt) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("if (%s) then (%s)\n", n.Condition, n.PositiveText)
	if len(n.PositiveStmt) != 0 {
		_ = n.PositiveStmt.Render(w)
	}

	if len(n.NegativeStmt) != 0 {
		w.Printf("else (%s)\n", n.NegativeText)
		_ = n.NegativeStmt.Render(w)
	}

	w.Print("endif\n")

	return w.Err
}

func (n *ActIfStmt) acStmt() {}

type ActGotoLabel struct {
	Name string
}

func (n *ActGotoLabel) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("label %s\n", n.Name)
	return w.Err
}

func (n *ActGotoLabel) acStmt() {}

type ActGoto struct {
	Name string
}

func (n *ActGoto) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Printf("goto %s\n", n.Name)
	return w.Err
}

func (n *ActGoto) acStmt() {}

// ============
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

func (a *Activity) AddStmt(stmt *Stmt) *Activity {
	a.Stmts = append(a.Stmts, stmt)
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
	Start *StartStmt
	While *WhileStmt
	Stop  *StopStmt

	Block []*Stmt

	Note      *ActivityNote
	Swimlane  *Swimlane
	SplitStmt *ActSplitStmt
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

	if n.Swimlane != nil {
		if err := n.Swimlane.Render(wr); err != nil {
			return err
		}
	}

	if n.While != nil {
		if err := n.While.Render(wr); err != nil {
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

	if n.SplitStmt != nil {
		if err := n.SplitStmt.Render(wr); err != nil {
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
