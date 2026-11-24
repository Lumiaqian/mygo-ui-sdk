# Composable Behaviors (`mygo-ui-sdk/behaviors`)

Every behavior operates in MyGo's immediate-mode frame pass using public `ui` APIs only.

## 1. `behaviors.TrackDrag` and `behaviors.SwitchTrack`

Use `TrackDrag` for any 1D draggable thumb or lens that must distinguish a tap from a drag and schedule frames during drag so release animations settle:

```go
res := behaviors.SwitchTrack(c.Key("auto-sync"), &app.sync, behaviors.SwitchTrackOptions{
	Label:    "Auto Sync",
	Min:      3,
	Max:      20,
	Disabled: disabled,
})
// Style res.Switch and paint the knob using res.Target / res.Progress / res.Active.
```

## 2. `behaviors.SegmentedTrack` and `behaviors.SegmentedGroup`

- Use `SegmentedTrack` when a single glass lens drags continuously across segments:
  - Handles pointer drag + snap, tap-to-segment, `Enter`/`Space` cycle, and focused `Left`/`Right` (RTL-aware), `Home`, `End`.
  - Sets stable `.Label(opts.Label)` and dynamic `.Value(labels[*selected])`.
- Use `SegmentedGroup` when each segment is its own focusable button on `ui.SegmentedBase`.

## 3. `behaviors.Toolbar`

Creates a `RoleToolbar` + `FocusGroup(orient)` container with zero visual chrome:

```go
tb := behaviors.Toolbar(c.Key("tools"), behaviors.ToolbarOptions{
	Label: "Annotation Tools",
})
tb.Root.Gap(4).Children(func() {
	if tb.Toggle(c.Key("rect"), "Rect", active == ToolRect, behaviors.ToolbarItemOptions{
		Tooltip:     "Rectangle (1)",
		ShortcutKey: ui.Key1,
	}).Size(28, 28).Clicked() {
		active = ToolRect
	}
	if tb.RadioItem(c.Key("red"), "Red", color == Red, behaviors.ToolbarItemOptions{}).Size(18, 18).Clicked() {
		color = Red
	}
})
```

## 4. `behaviors.PlaceFloating`

Pure geometry helper for floating toolbars, loupes, and popovers:
- Tries `PreferredSide`, flips to the opposite side if it overflows `viewport`, optionally falls back inside `anchor` (`InsideFallback: true`), clamps to `Margin`, and mirrors `AlignStart`/`AlignEnd` under `ui.RightToLeft`.

## 5. `behaviors.NavList`

Headless `RoleTabList` + `RoleTab` switcher for sidebar or horizontal navigation bars:
- Handles `Up`/`Down` (vertical, default) or `Left`/`Right` (horizontal, RTL-aware) plus `Home`/`End`, suppressed automatically when `c.FocusTakesText()` is true.
- Clamps at ends by default (`Wrap: true` to loop). Style `nav.Root` before `nav.Root.Children(...)` and build each entry with `nav.Item(c.Key(id), i, label)`.

## 6. `behaviors.RadioGroup`

Headless `RoleRadioGroup` + `FocusGroup(orient)` container for custom swatch rows or card choice grids:
- Style `rg.Root` before `rg.Root.Children(...)` and build choices with `rg.Item(c.Key(id), label, active, behaviors.RadioItemOptions{...})`.
- Inherits single Tab stop (landing on the checked radio), arrow-key roving focus, and `Home`/`End`.

## 7. `behaviors.ShouldDismiss`

Checks `Escape` (optionally suppressed while `c.FocusTakesText()`) and `surface.PressedOutside()`, returning `(dismissed, reason)`.
