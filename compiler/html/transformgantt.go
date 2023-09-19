package html

import (
	"github.com/worldiety/dddl/parser"
	"github.com/worldiety/dddl/plantuml"
	"github.com/worldiety/dddl/resolver"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"html/template"
	"log"
	"math"
)

func newProjectPlan(r *resolver.Resolver, model PreviewModel) *ProjectPlan {
	pp := &ProjectPlan{}
	tmp := map[string]ProjectTask{}
	parser.MustWalk(r.Workspace(), func(n parser.Node) error {
		if tDef, ok := n.(*parser.TypeDefinition); ok {
			wp := parser.FindAnnotation[*parser.WorkPackageAnnotation](tDef)
			if wp == nil {
				return nil
			}

			if wp.Name == "" {
				return nil
			}

			pt := tmp[wp.Name]
			pt.Name = wp.Name
			pt.Duration += wp.Duration
			pt.Refs = append(pt.Refs, resolver.NewQualifiedNameFromNamedType(tDef.Type).String())
			pt.Requires = append(wp.Requires, wp.Requires...)

			tmp[wp.Name] = pt
		}
		return nil
	})

	if len(tmp) == 0 {
		return nil
	}

	tasks := maps.Values(tmp)
	slices.SortFunc(tasks, func(a, b ProjectTask) bool {
		return a.Name < b.Name
	})

	for _, task := range tasks {
		task := task
		slices.Sort(task.Refs)
		task.Refs = slices.Compact(task.Refs)

		slices.Sort(task.Requires)
		task.Requires = slices.Compact(task.Requires)

		var tmpReq []string
		for _, require := range task.Requires {
			if _, ok := tmp[require]; ok {
				tmpReq = append(tmpReq, require)
			} else {
				log.Printf("gantt: required task does not exist: %s\n", require)
			}
		}
		task.Requires = tmpReq

		pp.Tasks = append(pp.Tasks, &task)
	}

	svg, err := plantuml.RenderLocal("svg", RenderGantt(pp))
	if err != nil {
		slog.Error("failed to convert choice to puml", slog.Any("err", err))
	}
	pp.GanttChartSVG = template.HTML(svg)

	return pp
}

func RenderGantt(pp *ProjectPlan) *plantuml.Diagram {
	diag := plantuml.NewDiagram()
	diag.BackgroundColor = "#00000000"
	chart := &plantuml.GanttChart{}
	for _, task := range pp.Tasks {
		t := &plantuml.GanttTask{
			Name:         task.Name,
			DurationDays: int(math.Ceil(task.Duration.Hours() / 8)),
			DependsOn:    task.Requires,
		}

		chart.Tasks = append(chart.Tasks, t)
	}
	diag.Add(chart)

	return diag
}
