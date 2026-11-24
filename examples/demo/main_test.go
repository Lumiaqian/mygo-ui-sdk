package demo_test

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/Lumiaqian/mygo-ui-sdk/behaviors"
	"github.com/Lumiaqian/mygo-ui-sdk/uitest"
)

func TestInspectorDemoView(t *testing.T) {
	mode := 0
	autoSave := true
	activeTool := "Rect"
	swatch := "Red"

	view := func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(12).Children(func() {
			seg := behaviors.SegmentedTrack(c.Key("mode"), &mode, behaviors.SegmentedTrackOptions{
				Label:        "Capture Mode",
				Labels:       []string{"Region", "Window", "Scroll"},
				SegmentWidth: 72,
				Padding:      3,
			})
			seg.Track.Size(222, 30).Children(func() {
				for _, lb := range []string{"Region", "Window", "Scroll"} {
					ui.Text(c, lb)
				}
			})

			sw := behaviors.SwitchTrack(c.Key("autosave"), &autoSave, behaviors.SwitchTrackOptions{
				Label: "Auto Save",
				Min:   3,
				Max:   20,
			})
			sw.Switch.Size(50, 28)

			pos := behaviors.PlaceFloating(
				image.Rect(0, 0, 640, 480),
				image.Rect(100, 100, 420, 260),
				image.Pt(240, 36),
				behaviors.PlacementOptions{
					PreferredSide:  behaviors.SideBottom,
					Align:          behaviors.AlignEnd,
					Gap:            8,
					Direction:      c.Direction(),
					InsideFallback: true,
				},
			)

			tb := behaviors.Toolbar(c.Key("toolbar"), behaviors.ToolbarOptions{
				Label: "Tools",
			})
			tb.Root.Absolute().Left(float32(pos.Point.X)).Top(float32(pos.Point.Y)).Gap(6).Children(func() {
				for _, tool := range []string{"Rect", "Arrow"} {
					if tb.Toggle(c.Key("tool-"+tool), tool, activeTool == tool, behaviors.ToolbarItemOptions{
						Tooltip: tool,
					}).Size(48, 28).Clicked() {
						activeTool = tool
					}
				}
				for _, col := range []string{"Red", "Blue"} {
					if tb.RadioItem(c.Key("swatch-"+col), col, swatch == col, behaviors.ToolbarItemOptions{}).Size(20, 20).Clicked() {
						swatch = col
					}
				}
			})
		})
	}

	tester := ui.NewTester(view, 640, 480)
	opts := uitest.AuditOptions{
		RequiredLabels:   []string{"Capture Mode", "Auto Save", "Tools", "Rect", "Arrow", "Red", "Blue"},
		RequireFocusable: true,
	}
	uitest.AssertAgentReady(t, tester, opts)
	uitest.AssertThemeParity(t, tester, opts)
	uitest.AssertDirectionParity(t, tester, opts)
}
