# Testing and Accessibility (`ui.Tester` + `uitest`)

Run tests windowlessly with `ui.NewTester` inside standard `go test`.

## 1. Minimum test for every new or refactored view

```go
tester := ui.NewTester(viewFn, 800, 600)

opts := uitest.AuditOptions{
	RequiredLabels:   []string{"Save", "Cancel"},
	RequireFocusable: true,
}
uitest.AssertAgentReady(t, tester, opts)
uitest.AssertThemeParity(t, tester, opts)
uitest.AssertDirectionParity(t, tester, opts)
```

## 2. What `uitest.Audit` catches

- `FindingUnnamedControl`: any interactive or focusable node with empty `Name`.
- `FindingDuplicateName`: multiple visible interactive controls sharing the same `Name` (breaks unambiguous agent targeting).
- `FindingMissingRequired`: a label listed in `RequiredLabels` or `RequiredInteractive` is absent from the visible accessibility tree.
- `FindingNonInteractiveRole`: a label listed in `RequiredInteractive` exists only as a passive container (e.g., raw `ui.Row`/`ui.Box` with `.Clicked()` instead of `ui.ButtonBase` or `behaviors`).
- `FindingNotFocusable`: an enabled interactive control is missing `Focusable` (when `RequireFocusable: true`).

## 3. Keyboard navigation verification

- `uitest.AssertKeyboardTabOrder(t, tester, []string{"First", "Second"})` checks sequential `Tab` order.
- `uitest.AssertRovingFocus(t, tester, ui.Horizontal, []string{"A", "B", "C"})` checks arrow keys, wrap-around, `Home`, and `End` inside a `FocusGroup`.
- Use `uitest.Snapshot(tester)` in failure logs to inspect every visible accessible node's `Role`, `Name`, `Value`, `Checked`, `Focusable`, and `Bounds`.

## 4. Editable fields, shortcuts and clipboard side effects

Give the editable element its own accessible name, for example
`ui.TextInput(c, &draft).Label("Annotation")`. A populated value can still
belong to an unnamed field. Audit the frame with the editor open.

Guard window shortcuts with `!c.FocusTakesText()` before calling
`c.Shortcut(...)`. Exercise the same key with the field unfocused and
focused. Use `tester.TypeKey(0, ui.Key1, "1")` to check both the resulting
text and the absence of a tool change.

For view-local text copies, use `c.Services().WriteClipboard(text)` and
assert `tester.Clipboard()`. For application code using native clipboard
APIs, inject a copy callback and capture its argument in the test.
`ui.NewTester` supplies a headless host. An off-main-thread
`mygo.RunOnMain` waits for the native event loop, so exercise that path in
a native integration test with a running loop.
