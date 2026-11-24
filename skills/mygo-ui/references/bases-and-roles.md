# Choosing Bases, Behaviors, and Roles

## 1. Decision order

1. **MyGo built-in styled widget** (`ui.Button`, `ui.Checkbox`, `ui.Select`, `ui.Toolbar`) — when default theme styling is enough.
2. **MyGo headless `*Base`** — when you want custom styling with standard desktop interaction:
   - `ui.ButtonBase` — clickable action button (pointer + Enter/Space + focus)
   - `ui.CheckboxBase` / `ui.CheckboxStateBase` / `ui.SwitchBase` / `ui.ToggleBase` — boolean or tri-state toggle
   - `ui.RadioBase` — mutually exclusive choice
   - `ui.SegmentedBase` — discrete segmented radio group with arrow keys
   - `ui.TabsBase`, `ui.AccordionBase`, `ui.CollapsibleBase`, `ui.TreeBase`
   - `ui.SelectBase`, `ui.ComboboxBase`, `ui.MenuBase`, `ui.ContextMenuBase`
   - `ui.PopoverBase`, `ui.DialogBase`, `ui.TooltipBase`, `ui.ToastBase`
   - `ui.SliderBase`, `ui.TextInputBase`, `ui.TextAreaBase`
3. **`behaviors.*` from `mygo-ui-sdk`** — when you need custom interaction physics or unstyled composition beyond a single `*Base`:
   - Continuous draggable glass lens segmented control → `behaviors.SegmentedTrack`
   - Draggable custom switch knob → `behaviors.SwitchTrack`
   - Discrete custom segmented control with pre-labeled segments → `behaviors.SegmentedGroup`
   - Floating / glass toolbar with roving focus, toggle buttons, and color/stroke swatches → `behaviors.Toolbar`
   - Viewport-safe anchor placement with flip and RTL mirroring → `behaviors.PlaceFloating`
   - Light-dismiss on Escape / outside press → `behaviors.ShouldDismiss`

## 2. Common anti-patterns to avoid

- **Clickable `ui.Box` / `ui.Row`**: Clicking a `ui.Box` works with a mouse, but it has no `Focusable()`, no Enter/Space activation, and may disappear from the accessibility tree. Always start from `ui.ButtonBase` or `tb.Button` / `tb.Toggle` / `tb.RadioItem`.
- **Icon-only or color-dot buttons without `Label`**: Always call `.Label("红")` / `.Label("删除")` in addition to `.Tooltip(...)`. Tooltips populate `Description`, not `Name`.
- **Encoding dynamic value into `Label` instead of `Value`**: Prefer `.Label("画质").Value("高")` over `.Label("画质：高")` when possible so agents and tests can locate the control by stable name `"画质"`.
- **Forgetting `Checked(active)` on toggle/radio controls**: Visual highlight alone is invisible to assistive technology and agents. Set `.Role(ui.RoleToggleButton).Checked(active)` or `.Role(ui.RoleRadio).Checked(active)`.
