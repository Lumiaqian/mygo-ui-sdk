package uitest

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// AssertThemeParity verifies Audit passes in both Light and Dark themes.
func AssertThemeParity(tb testing.TB, t *ui.Tester, opts AuditOptions) {
	tb.Helper()
	for _, tc := range []struct {
		name string
		dark bool
	}{
		{"light", false},
		{"dark", true},
	} {
		t.SetDark(tc.dark)
		t.Frame()
		findings := Audit(t, opts)
		if len(findings) > 0 {
			tb.Errorf("theme=%s failed agent-readiness audit:\n%s", tc.name, formatFindings(findings))
		}
	}
}

// AssertDirectionParity verifies Audit passes in both LTR and RTL directions
// and that the set of visible interactive control labels stays consistent.
func AssertDirectionParity(tb testing.TB, t *ui.Tester, opts AuditOptions) {
	tb.Helper()

	t.SetDirection(ui.LeftToRight)
	t.Frame()
	if findings := Audit(t, opts); len(findings) > 0 {
		tb.Errorf("direction=LTR failed agent-readiness audit:\n%s", formatFindings(findings))
	}
	ltrLabels := interactiveLabels(t, opts.IncludeOffscreen)

	t.SetDirection(ui.RightToLeft)
	t.Frame()
	if findings := Audit(t, opts); len(findings) > 0 {
		tb.Errorf("direction=RTL failed agent-readiness audit:\n%s", formatFindings(findings))
	}
	rtlLabels := interactiveLabels(t, opts.IncludeOffscreen)

	t.SetDirection(ui.LeftToRight)
	t.Frame()

	if !slices.Equal(ltrLabels, rtlLabels) {
		tb.Errorf("interactive labels diverged between LTR %v and RTL %v", ltrLabels, rtlLabels)
	}
}

// AssertKeyboardTabOrder presses Tab once per expected label and checks focus moves in order.
func AssertKeyboardTabOrder(tb testing.TB, t *ui.Tester, expectedNames []string) {
	tb.Helper()
	for i, want := range expectedNames {
		t.Key(0, ui.KeyTab)
		if !t.Focused(want) {
			tb.Fatalf("tab step %d: expected %q to be focused; current focus=%q", i, want, focusedName(t))
		}
	}
}

// AssertRovingFocus verifies arrow-key roving focus across names inside a FocusGroup,
// including Home, End, and wrap-around. The first item must already be focused.
func AssertRovingFocus(tb testing.TB, t *ui.Tester, orient ui.Orientation, names []string) {
	tb.Helper()
	if len(names) < 2 {
		tb.Fatalf("AssertRovingFocus requires at least 2 items, got %d", len(names))
	}
	nextKey, prevKey := ui.KeyRight, ui.KeyLeft
	if orient == ui.Vertical {
		nextKey, prevKey = ui.KeyDown, ui.KeyUp
	}

	if !t.Focused(names[0]) {
		tb.Fatalf("expected initial focus on %q; current focus=%q", names[0], focusedName(t))
	}
	for i := 1; i < len(names); i++ {
		t.Key(0, nextKey)
		if !t.Focused(names[i]) {
			tb.Fatalf("step %d forward: expected %q focused; current focus=%q", i, names[i], focusedName(t))
		}
	}
	// Wrap around to first item.
	t.Key(0, nextKey)
	if !t.Focused(names[0]) {
		tb.Fatalf("wrap forward: expected %q focused; current focus=%q", names[0], focusedName(t))
	}
	// Wrap backward to last item.
	t.Key(0, prevKey)
	last := names[len(names)-1]
	if !t.Focused(last) {
		tb.Fatalf("wrap backward: expected %q focused; current focus=%q", last, focusedName(t))
	}
	// Home jumps to first, End jumps to last.
	t.Key(0, ui.KeyHome)
	if !t.Focused(names[0]) {
		tb.Fatalf("Home: expected %q focused; current focus=%q", names[0], focusedName(t))
	}
	t.Key(0, ui.KeyEnd)
	if !t.Focused(last) {
		tb.Fatalf("End: expected %q focused; current focus=%q", last, focusedName(t))
	}
}

func interactiveLabels(t *ui.Tester, includeOffscreen bool) []string {
	var out []string
	for _, n := range t.AccessAll() {
		if n.Offscreen && !includeOffscreen {
			continue
		}
		if InteractiveRole(n.Role) {
			out = append(out, strings.TrimSpace(n.Name))
		}
	}
	slices.Sort(out)
	return out
}

func focusedName(t *ui.Tester) string {
	for _, n := range t.AccessAll() {
		if n.Focused {
			return n.Name
		}
	}
	return ""
}

func formatFindings(findings []Finding) string {
	var b strings.Builder
	for _, f := range findings {
		fmt.Fprintf(&b, "  - %s\n", f.String())
	}
	return b.String()
}
