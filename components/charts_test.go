package components

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/twaananen/boulderlog/services"
)

func TestChartEmptyStates(t *testing.T) {
	for name, component := range map[string]templ.Component{
		"counts":      GradeCountsChart(nil, nil, "weekly", "2026-10-25", true),
		"progression": DifficultyProgressionChart(nil, nil, "week"),
	} {
		t.Run(name, func(t *testing.T) {
			var html bytes.Buffer
			if err := component.Render(context.Background(), &html); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(html.String(), "No climbs logged for this period.") || strings.Contains(html.String(), "<canvas") {
				t.Fatalf("expected an empty state without a canvas: %s", html.String())
			}
		})
	}
}

func TestChartsRenderDataWithoutInlineScripts(t *testing.T) {
	for name, component := range map[string]templ.Component{
		"counts": GradeCountsChart([]string{"6A"}, map[string][]int{"Topped": {2}}, "all", "2026-10-25", false),
		"progression": DifficultyProgressionChart(map[string][]services.DifficultyDataPoint{
			"6A": {{Value: 3, Label: "W1"}},
		}, []string{"W1"}, "week"),
	} {
		t.Run(name, func(t *testing.T) {
			var html bytes.Buffer
			if err := component.Render(context.Background(), &html); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(html.String(), `data-chart="`) || strings.Contains(html.String(), "<script") {
				t.Fatalf("expected chart data without an inline initializer: %s", html.String())
			}
		})
	}
}
