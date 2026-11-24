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
