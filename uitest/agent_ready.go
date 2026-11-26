package uitest

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
)

// FindingKind classifies an accessibility or agent-readiness issue found by Audit.
type FindingKind string

const (
	FindingUnnamedControl     FindingKind = "unnamed_control"
	FindingDuplicateName      FindingKind = "duplicate_name"
	FindingMissingRequired    FindingKind = "missing_required"
	FindingNotFocusable       FindingKind = "not_focusable"
	FindingNonInteractiveRole FindingKind = "non_interactive_role"
)

// Finding describes one issue detected in the last rendered accessibility tree.
type Finding struct {
	Kind   FindingKind
	Role   ui.Role
	Name   string
	Detail string
}

func (f Finding) String() string {
	if f.Name != "" {
		return fmt.Sprintf("[%s] role=%d name=%q: %s", f.Kind, f.Role, f.Name, f.Detail)
	}
	return fmt.Sprintf("[%s] role=%d: %s", f.Kind, f.Role, f.Detail)
}

// AuditOptions configures Audit and AssertAgentReady.
type AuditOptions struct {
	// RequiredLabels lists accessible names that must appear in the visible tree.
	RequiredLabels []string
	// RequiredInteractive lists accessible names that must appear as visible,
	// interactive roles (and focusable when not disabled and RequireFocusable is true).
	// This catches raw ui.Box/ui.Row elements that call Clicked() without ui.ButtonBase.
	RequiredInteractive []string
	// AllowDuplicateNames exempts specific names from duplicate interactive label checks.
	AllowDuplicateNames []string
	// IgnoreUnnamedRoles exempts specific roles from unnamed-control checks.
	IgnoreUnnamedRoles []ui.Role
	// RequireFocusable checks that enabled, visible interactive controls are keyboard focusable.
	RequireFocusable bool
	// IgnoreOffscreen skips offscreen nodes (enabled by default when IncludeOffscreen is false).
	IncludeOffscreen bool
}

// SnapshotEntry is a compact representation of an accessible node for agent diagnostics.
type SnapshotEntry struct {
	Role         ui.Role
	Name         string
	Value        string
	Description  string
	KeyShortcuts string
	Disabled     bool
	Focusable    bool
	Focused      bool
	Checked      bool
	Selected     bool
	Expanded     bool
	Offscreen    bool
	Bounds       ui.Rect
}

// Snapshot returns a compact view of the accessibility tree from the last frame.
func Snapshot(t *ui.Tester) []SnapshotEntry {
	nodes := t.AccessAll()
	out := make([]SnapshotEntry, 0, len(nodes))
	for _, n := range nodes {
		out = append(out, SnapshotEntry{
			Role:         n.Role,
			Name:         n.Name,
			Value:        n.Value,
			Description:  n.Description,
			KeyShortcuts: n.KeyShortcuts,
			Disabled:     n.Disabled,
			Focusable:    n.Focusable,
			Focused:      n.Focused,
			Checked:      n.Checked,
			Selected:     n.Selected,
			Expanded:     n.Expanded,
			Offscreen:    n.Offscreen,
			Bounds:       n.Bounds,
		})
	}
	return out
}

// FindByRole returns all AccessInfo nodes matching role in the last frame.
func FindByRole(t *ui.Tester, role ui.Role) []ui.AccessInfo {
	var out []ui.AccessInfo
	for _, n := range t.AccessAll() {
		if n.Role == role {
			out = append(out, n)
		}
	}
	return out
}

// InteractiveRole reports whether r is a control role an agent or keyboard user acts on directly.
func InteractiveRole(r ui.Role) bool {
	switch r {
	case ui.RoleButton,
		ui.RoleLink,
		ui.RoleCheckBox,
		ui.RoleRadio,
		ui.RoleSwitch,
		ui.RoleSlider,
		ui.RoleTextField,
		ui.RolePopUpButton,
		ui.RoleTab,
		ui.RoleMenuButton,
		ui.RoleToggleButton,
		ui.RoleComboBox,
		ui.RoleDisclosure,
		ui.RoleStepper,
		ui.RoleColorWell,
		ui.RoleMenuItem,
		ui.RoleMenuItemCheckBox,
		ui.RoleMenuItemRadio:
		return true
	default:
		return false
	}
}

// Audit inspects t.AccessAll() and returns accessibility / agent-readiness findings.
func Audit(t *ui.Tester, opts AuditOptions) []Finding {
	nodes := t.AccessAll()
	var findings []Finding

	seenNames := map[string]struct{}{}
	nodesByName := map[string][]ui.AccessInfo{}
	interactiveCounts := map[string]int{}
	interactiveRoles := map[string]ui.Role{}

	for _, n := range nodes {
		if n.Offscreen && !opts.IncludeOffscreen {
			continue
		}
		name := strings.TrimSpace(n.Name)
		if name != "" {
			seenNames[name] = struct{}{}
			nodesByName[name] = append(nodesByName[name], n)
		}

		interactive := InteractiveRole(n.Role)
		if !interactive && !n.Focusable {
			continue
		}

		if name == "" && !slices.Contains(opts.IgnoreUnnamedRoles, n.Role) {
			hint := "set Label(...) or inner text"
			if n.Role == ui.RoleTextField {
				hint = "set Label(...) on the text field itself"
			}
			findings = append(findings, Finding{
				Kind:   FindingUnnamedControl,
				Role:   n.Role,
				Detail: fmt.Sprintf("focusable=%v bounds=%v has no accessible name; %s", n.Focusable, n.Bounds, hint),
			})
		}

		if interactive && !n.Disabled && opts.RequireFocusable && !n.Focusable {
			findings = append(findings, Finding{
				Kind:   FindingNotFocusable,
				Role:   n.Role,
				Name:   name,
				Detail: "enabled interactive control is not keyboard focusable",
			})
		}

		if interactive && name != "" {
			interactiveCounts[name]++
			interactiveRoles[name] = n.Role
		}
	}

	dupNames := make([]string, 0)
	for name, count := range interactiveCounts {
		if count > 1 && !slices.Contains(opts.AllowDuplicateNames, name) {
			dupNames = append(dupNames, name)
		}
	}
	slices.Sort(dupNames)
	for _, name := range dupNames {
		findings = append(findings, Finding{
			Kind:   FindingDuplicateName,
			Role:   interactiveRoles[name],
			Name:   name,
			Detail: fmt.Sprintf("appears %d times among interactive controls; disambiguate labels so agents can target it unambiguously", interactiveCounts[name]),
		})
	}

	for _, req := range opts.RequiredLabels {
		if _, ok := seenNames[req]; !ok {
			findings = append(findings, Finding{
				Kind:   FindingMissingRequired,
				Name:   req,
				Detail: "required accessible label not found in visible accessibility tree",
			})
		}
	}

	for _, req := range opts.RequiredInteractive {
		matches := nodesByName[req]
		if len(matches) == 0 {
			findings = append(findings, Finding{
				Kind:   FindingMissingRequired,
				Name:   req,
				Detail: "required interactive control not found in visible accessibility tree",
			})
			continue
		}
		var hasInteractive bool
		var hasFocusable bool
		var firstRole ui.Role
		for _, m := range matches {
			if firstRole == 0 {
				firstRole = m.Role
			}
			if InteractiveRole(m.Role) {
				hasInteractive = true
				if m.Focusable || m.Disabled {
					hasFocusable = true
				}
			}
		}
		if !hasInteractive {
			findings = append(findings, Finding{
				Kind:   FindingNonInteractiveRole,
				Role:   firstRole,
				Name:   req,
				Detail: "found in accessibility tree, but not as an interactive control role (likely raw ui.Box/ui.Row with Clicked instead of ui.ButtonBase or behaviors)",
			})
		} else if opts.RequireFocusable && !hasFocusable {
			findings = append(findings, Finding{
				Kind:   FindingNotFocusable,
				Role:   firstRole,
				Name:   req,
				Detail: "interactive control exists but is not keyboard focusable",
			})
		}
	}

	return findings
}

// AssertAgentReady fails tb with a formatted summary if Audit reports any findings.
func AssertAgentReady(tb testing.TB, t *ui.Tester, opts AuditOptions) {
	tb.Helper()
	findings := Audit(t, opts)
	if len(findings) == 0 {
		return
	}
	var b strings.Builder
	fmt.Fprintf(&b, "agent-readiness audit found %d issue(s):\n", len(findings))
	for _, f := range findings {
		fmt.Fprintf(&b, "  - %s\n", f.String())
	}
	tb.Error(b.String())
}
