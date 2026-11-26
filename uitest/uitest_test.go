package uitest_test

import (
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/Lumiaqian/mygo-ui-sdk/uitest"
)

func TestAuditDetectsUnnamedDuplicateAndMissing(t *testing.T) {
	tester := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			// Unnamed focusable button.
			ui.ButtonBase(c.Key("unnamed")).Size(32, 32)
			// Duplicate interactive label.
			ui.Button(c.Key("dup-1"), "Save")
			ui.Button(c.Key("dup-2"), "Save")
		})
	}, 400, 300)

	findings := uitest.Audit(tester, uitest.AuditOptions{
		RequiredLabels:   []string{"Save", "Cancel"},
		RequireFocusable: true,
	})

	wantKinds := map[uitest.FindingKind]bool{
		uitest.FindingUnnamedControl:  false,
		uitest.FindingDuplicateName:   false,
		uitest.FindingMissingRequired: false,
	}
	for _, f := range findings {
		if _, ok := wantKinds[f.Kind]; ok {
			wantKinds[f.Kind] = true
		}
	}
	for kind, found := range wantKinds {
		if !found {
			t.Errorf("expected finding kind %s, got findings=%v", kind, findings)
		}
	}
}

func TestAssertAgentReadyAndParityPass(t *testing.T) {
	var synced bool
	tester := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			ui.Checkbox(c.Key("sync"), &synced, "Auto Sync")
			ui.Button(c.Key("apply"), "Apply")
		})
	}, 400, 300)

	opts := uitest.AuditOptions{
		RequiredLabels:   []string{"Auto Sync", "Apply"},
		RequireFocusable: true,
	}
	uitest.AssertAgentReady(t, tester, opts)
	uitest.AssertThemeParity(t, tester, opts)
	uitest.AssertDirectionParity(t, tester, opts)
	uitest.AssertKeyboardTabOrder(t, tester, []string{"Auto Sync", "Apply"})
}

func TestAuditDetectsNonInteractiveBoxClick(t *testing.T) {
	tester := ui.NewTester(func(c *ui.Context) {
		box := ui.Row(c.Key("fake-btn")).Label("Fake Button").Size(80, 28)
		_ = box.Clicked()
	}, 200, 100)

	findings := uitest.Audit(tester, uitest.AuditOptions{
		RequiredInteractive: []string{"Fake Button"},
		RequireFocusable:    true,
	})
	if len(findings) != 1 || findings[0].Kind != uitest.FindingNonInteractiveRole {
		t.Fatalf("expected FindingNonInteractiveRole, got %v", findings)
	}
}

func TestAssertRovingFocus(t *testing.T) {
	tester := ui.NewTester(func(c *ui.Context) {
		ui.Row(c.Key("tb")).Role(ui.RoleToolbar).Label("Actions").FocusGroup(ui.Horizontal).Gap(4).Children(func() {
			ui.Button(c.Key("a"), "Cut")
			ui.Button(c.Key("b"), "Copy")
			ui.Button(c.Key("c"), "Paste")
		})
	}, 400, 100)

	tester.Key(0, ui.KeyTab)
	uitest.AssertRovingFocus(t, tester, ui.Horizontal, []string{"Cut", "Copy", "Paste"})
}

func TestUnnamedTextField(t *testing.T) {
	value := "draft"
	label := ""
	tt := ui.NewTester(func(c *ui.Context) {
		ui.TextInput(c, &value).Label(label)
	}, 300, 100)
	findings := uitest.Audit(tt, uitest.AuditOptions{})
	if len(findings) != 1 || findings[0].Kind != uitest.FindingUnnamedControl ||
		findings[0].Role != ui.RoleTextField || !strings.Contains(findings[0].Detail, "Label(...) on the text field") {
		t.Fatalf("unexpected text field diagnostics: %v", findings)
	}
	label = "Annotation"
	tt.Frame()
	uitest.AssertAgentReady(t, tt, uitest.AuditOptions{RequiredInteractive: []string{label}, RequireFocusable: true})
}

func TestClipboardShortcutWhileEditing(t *testing.T) {
	var value string
	tt := ui.NewTester(func(c *ui.Context) {
		ui.TextInput(c, &value).Label("Annotation")
		if !c.FocusTakesText() && c.Shortcut(0, ui.Key1) {
			c.Services().WriteClipboard("tool 1")
		}
	}, 300, 100)
	tt.Key(0, ui.Key1)
	if got := tt.Clipboard(); got != "tool 1" {
		t.Fatalf("shortcut clipboard = %q", got)
	}
	tt.SetClipboard("original")
	if err := tt.Click("Annotation"); err != nil {
		t.Fatal(err)
	}
	tt.TypeKey(0, ui.Key1, "1")
	if got := tt.Clipboard(); got != "original" || value != "1" {
		t.Fatalf("while editing: clipboard=%q input=%q", got, value)
	}
}
