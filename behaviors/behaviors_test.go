package behaviors_test

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/Lumiaqian/mygo-ui-sdk/behaviors"
	"github.com/Lumiaqian/mygo-ui-sdk/uitest"
)

func TestSegmentedTrackKeyboardAndRTL(t *testing.T) {
	selected := 0
	labels := []string{"Day", "Week", "Month"}

	tester := ui.NewTester(func(c *ui.Context) {
		res := behaviors.SegmentedTrack(c.Key("period"), &selected, behaviors.SegmentedTrackOptions{
			Label:        "Period",
			Labels:       labels,
			SegmentWidth: 80,
			Padding:      4,
		})
		res.Track.Size(248, 32).Children(func() {
			for _, lb := range labels {
				ui.Text(c, lb)
			}
		})
	}, 400, 120)

	uitest.AssertAgentReady(t, tester, uitest.AuditOptions{
		RequiredLabels:   []string{"Period"},
		RequireFocusable: true,
	})

	// Tab to track and navigate with arrows, End, Home.
	tester.Key(0, ui.KeyTab)
	if !tester.Focused("Period") {
		t.Fatalf("expected Period focused")
	}
	tester.Key(0, ui.KeyRight)
	if selected != 1 {
		t.Fatalf("KeyRight: got %d, want 1", selected)
	}
	tester.Key(0, ui.KeyEnd)
	if selected != 2 {
		t.Fatalf("KeyEnd: got %d, want 2", selected)
	}
	tester.Key(0, ui.KeyHome)
	if selected != 0 {
		t.Fatalf("KeyHome: got %d, want 0", selected)
	}

	// Verify RTL reverses Left/Right direction.
	tester.SetDirection(ui.RightToLeft)
	tester.Frame()
	tester.Key(0, ui.KeyLeft)
	if selected != 1 {
		t.Fatalf("RTL KeyLeft: got %d, want 1", selected)
	}
	tester.Key(0, ui.KeyRight)
	if selected != 0 {
		t.Fatalf("RTL KeyRight: got %d, want 0", selected)
	}

	// Verify accessible Value tracks selected label.
	info, ok := tester.Access("Period")
	if !ok || info.Value != "Day" {
		t.Fatalf("Access(Period): ok=%v value=%q, want Day", ok, info.Value)
	}
}

func TestSegmentedGroupAgentReadyAndKeyboard(t *testing.T) {
	selected := 0
	labels := []string{"2x", "2.5x", "3x"}

	tester := ui.NewTester(func(c *ui.Context) {
		seg := behaviors.SegmentedGroup(c.Key("zoom"), &selected, behaviors.SegmentedGroupOptions{
			Label:  "Zoom",
			Labels: labels,
		})
		seg.Track.Gap(4).Children(func() {
			for i, lb := range labels {
				seg.Segment(i).Padding(4, 8).Children(func() {
					ui.Text(c, lb)
				})
			}
		})
	}, 400, 100)

	uitest.AssertAgentReady(t, tester, uitest.AuditOptions{
		RequiredLabels:   []string{"Zoom", "2x", "2.5x", "3x"},
		RequireFocusable: true,
	})

	if err := tester.Click("3x"); err != nil {
		t.Fatalf("click 3x: %v", err)
	}
	if selected != 2 {
		t.Fatalf("after click 3x: got %d, want 2", selected)
	}
}

func TestSwitchTrackToggle(t *testing.T) {
	on := false
	tester := ui.NewTester(func(c *ui.Context) {
		res := behaviors.SwitchTrack(c.Key("sync"), &on, behaviors.SwitchTrackOptions{
			Label: "Auto Sync",
			Min:   3,
			Max:   20,
		})
		res.Switch.Size(50, 28)
	}, 200, 80)

	uitest.AssertAgentReady(t, tester, uitest.AuditOptions{
		RequiredLabels:   []string{"Auto Sync"},
		RequireFocusable: true,
	})
	if err := tester.Click("Auto Sync"); err != nil {
		t.Fatalf("click switch: %v", err)
	}
	if !on {
		t.Fatalf("expected switch on after click")
	}
}

func TestToolbarRovingFocusAndRoles(t *testing.T) {
	activeTool := "Rect"
	activeColor := "Red"

	tester := ui.NewTester(func(c *ui.Context) {
		tb := behaviors.Toolbar(c.Key("annotate"), behaviors.ToolbarOptions{
			Label: "Annotation Tools",
		})
		tb.Root.Gap(4).Children(func() {
			for _, name := range []string{"Rect", "Arrow"} {
				btn := tb.Toggle(c.Key(name), name, activeTool == name, behaviors.ToolbarItemOptions{
					Tooltip: name + " tool",
				}).Size(40, 28)
				if btn.Clicked() {
					activeTool = name
				}
			}
			for _, col := range []string{"Red", "Blue"} {
				sw := tb.RadioItem(c.Key("col-"+col), col, activeColor == col, behaviors.ToolbarItemOptions{}).Size(20, 20)
				if sw.Clicked() {
					activeColor = col
				}
			}
		})
	}, 400, 100)

	uitest.AssertAgentReady(t, tester, uitest.AuditOptions{
		RequiredLabels:   []string{"Annotation Tools", "Rect", "Arrow", "Red", "Blue"},
		RequireFocusable: true,
	})
	uitest.AssertThemeParity(t, tester, uitest.AuditOptions{RequireFocusable: true})
	uitest.AssertDirectionParity(t, tester, uitest.AuditOptions{RequireFocusable: true})

	// MyGo's FocusGroup prioritizes a checked RoleRadio on initial Tab entry;
	// pressing Home jumps to the first toolbar control, then roving traverses all items.
	tester.Key(0, ui.KeyTab)
	if !tester.Focused("Red") {
		t.Fatalf("expected FocusGroup to prioritize checked radio Red, got focus elsewhere")
	}
	tester.Key(0, ui.KeyHome)
	uitest.AssertRovingFocus(t, tester, ui.Horizontal, []string{"Rect", "Arrow", "Red", "Blue"})

	if err := tester.Click("Blue"); err != nil {
		t.Fatalf("click Blue: %v", err)
	}
	if activeColor != "Blue" {
		t.Fatalf("activeColor=%q, want Blue", activeColor)
	}
	blueInfo, ok := tester.Access("Blue")
	if !ok || blueInfo.Role != ui.RoleRadio || !blueInfo.Checked {
		t.Fatalf("Blue access info: ok=%v role=%v checked=%v", ok, blueInfo.Role, blueInfo.Checked)
	}
}

func TestPlaceFloatingFlipClampAndRTL(t *testing.T) {
	vp := image.Rect(0, 0, 800, 600)
	size := image.Pt(200, 40)

	// Fits below anchor, aligned to End.
	res := behaviors.PlaceFloating(vp, image.Rect(300, 100, 500, 200), size, behaviors.PlacementOptions{
		PreferredSide:  behaviors.SideBottom,
		Align:          behaviors.AlignEnd,
		Gap:            8,
		InsideFallback: true,
	})
	if res.Side != behaviors.SideBottom || res.Flipped || res.Point != image.Pt(300, 208) {
		t.Fatalf("below: got %+v", res)
	}

	// Overflows bottom -> flips above.
	res = behaviors.PlaceFloating(vp, image.Rect(300, 500, 500, 580), size, behaviors.PlacementOptions{
		PreferredSide:  behaviors.SideBottom,
		Align:          behaviors.AlignEnd,
		Gap:            8,
		InsideFallback: true,
	})
	if res.Side != behaviors.SideTop || !res.Flipped || res.Point.Y != 452 {
		t.Fatalf("flip top: got %+v", res)
	}

	// Neither below nor above fits -> falls back inside anchor.
	res = behaviors.PlaceFloating(vp, image.Rect(100, 10, 500, 590), size, behaviors.PlacementOptions{
		PreferredSide:  behaviors.SideBottom,
		Align:          behaviors.AlignEnd,
		Gap:            8,
		InsideFallback: true,
	})
	if res.Point.Y != 542 {
		t.Fatalf("inside fallback: got %+v, want Y=542", res)
	}

	// RTL flips AlignEnd -> AlignStart.
	resRTL := behaviors.PlaceFloating(vp, image.Rect(300, 100, 600, 200), size, behaviors.PlacementOptions{
		PreferredSide: behaviors.SideBottom,
		Align:         behaviors.AlignEnd,
		Gap:           8,
		Direction:     ui.RightToLeft,
	})
	if resRTL.Point.X != 300 {
		t.Fatalf("RTL AlignEnd->Start: got X=%d, want 300", resRTL.Point.X)
	}
}

func TestShouldDismissEscapeAndOutside(t *testing.T) {
	open := true
	var reason behaviors.DismissReason

	tester := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Fill().Children(func() {
			if open {
				panel := ui.Box(c.Key("panel")).Absolute().Top(40).Left(40).Size(120, 80)
				if ok, r := behaviors.ShouldDismiss(c, panel, behaviors.DismissOptions{}); ok {
					open = false
					reason = r
				}
			}
		})
	}, 400, 300)

	tester.Key(0, ui.KeyEscape)
	if open || reason != behaviors.DismissEscape {
		t.Fatalf("Escape dismiss: open=%v reason=%v", open, reason)
	}

	open = true
	reason = behaviors.DismissNone
	tester.Frame()
	tester.ClickAt(10, 10)
	if open || reason != behaviors.DismissOutsidePress {
		t.Fatalf("Outside dismiss: open=%v reason=%v", open, reason)
	}
}
