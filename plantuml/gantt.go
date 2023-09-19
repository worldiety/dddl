package plantuml

import (
	"io"
	"time"
)

type GanttTask struct {
	Name         string
	DurationDays int
	DependsOn    []string
}

func (g *GanttTask) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Print("[")
	w.Print(g.Name)
	w.Print("]")
	if g.DurationDays != 0 {
		w.Printf(" lasts %d days", g.DurationDays)
	} else {
		w.Printf(" lasts 1 days")
	}
	w.Print("\n")

	for _, otherTaskName := range g.DependsOn {
		w.Print("[")
		w.Print(g.Name)
		w.Print("] starts at [")
		w.Print(otherTaskName)
		w.Print("]'s end\n")
	}

	w.Print("\n")
	return w.Err
}

type GanttChart struct {
	StartAt time.Time
	Tasks   []*GanttTask
}

func (g *GanttChart) Render(wr io.Writer) error {
	w := strWriter{Writer: wr}
	w.Print("@startgantt\n")
	if !g.StartAt.IsZero() {
		w.Printf("project starts %s\n", g.StartAt.Format(time.DateOnly))
		w.Print("\n")
	}

	for _, task := range g.Tasks {
		_ = task.Render(w)
	}

	w.Print("@endgantt\n")

	return w.Err
}
