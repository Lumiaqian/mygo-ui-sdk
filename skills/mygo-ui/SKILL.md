---
name: mygo-ui
description: Builds accessible, keyboard-navigable, agent-ready Native GPU interfaces with github.com/egoist/mygo/ui and github.com/Lumiaqian/mygo-ui-sdk. Covers choosing ui.*Base vs composable behaviors, Context.Key and ui.Local lifecycle rules, naming and ARIA semantics, and automated ui.Tester / uitest conformance checks. Use whenever building or refactoring Go native UI views in a MyGo app.
---

# Build Native UI with MyGo + mygo-ui-sdk

`github.com/egoist/mygo/ui` is an immediate-mode GPU UI toolkit in Go (`func(c *ui.Context)`). It already ships 18+ headless `*Base` controls (`ButtonBase`, `SwitchBase`, `SegmentedBase`, `SelectBase`, `TabsBase`, `MenuBase`, `DialogBase`, etc.), a full accessibility tree, and windowless `ui.Tester`.

`mygo-ui-sdk` adds the missing layers for AI agents and custom-styled controls (such as Liquid Glass or floating canvas toolbars):
- `behaviors`: composable immediate-mode interaction behaviors (`TrackDrag`, `SwitchTrack`, `SegmentedTrack`, `SegmentedGroup`, `Toolbar`, `PlaceFloating`, `ShouldDismiss`).
- `uitest`: automated accessibility, keyboard, direction (LTR/RTL), theme, and agent-readiness audits over `ui.Tester`.

## Read first (progressive disclosure)

Read [`llms.txt`](llms.txt) first, then open only the reference file your task needs:
- Choosing a base widget or ARIA role → [`references/bases-and-roles.md`](references/bases-and-roles.md)
- Building custom draggable / floating / toolbar controls → [`references/behaviors.md`](references/behaviors.md)
- Managing frame state, `Context.Key`, `ui.Local`, and focus → [`references/state-and-keys.md`](references/state-and-keys.md)
- Writing `ui.Tester` and `uitest.AssertAgentReady` tests → [`references/testing-and-a11y.md`](references/testing-and-a11y.md)

## Non-negotiable rules

1. **Never hand-roll interactive `ui.Box` / `ui.Row` clicks.** If an element calls `.Clicked()`, it MUST be built on `ui.ButtonBase`, another `ui.*Base`, or a `behaviors.*` helper so it gets keyboard focus (`Focusable`), Enter/Space activation, and an accessible role. A raw `ui.Box(c)` with `.Clicked()` is invisible to keyboard users and agents.
2. **Name every interactive control.** Every button, toggle, swatch, input, slider, and toolbar must have a visible text child or an explicit `.Label(...)`. Avoid ambiguous duplicate labels across visible interactive controls.
3. **Prefer `ui.*Base` first, then `behaviors.*`.** Use `ui.SegmentedBase` / `ui.SwitchBase` / `ui.SelectBase` / `ui.MenuBase` whenever the interaction matches. Drop down to `behaviors.*` when building custom physics or chrome (e.g., draggable Liquid Glass lenses, floating toolbars, custom popovers).
4. **Always key dynamic or stateful controls.** Pass `c.Key(...)` before constructing stateful bases, lists, loops, or elements storing `ui.Local` / `AnimateTo` state.
5. **Respect `FocusGroup` radio priority.** Inside a `FocusGroup`, initial `Tab` lands on a checked `RoleRadio` item if present (else the last focused or first item); `Home`/`End` jump to ends.
6. **Keep `ui.Element` ephemeral.** Never store `ui.Element` or `*ui.Context` in app structs or background goroutines; capture `c.Services()` or `ui.Handle` instead.
7. **Always close the loop with `uitest`.** Every new or refactored view must have a `ui.NewTester` test running `uitest.AssertAgentReady` (with `RequireFocusable: true` and `RequiredLabels`), plus `AssertThemeParity` and `AssertDirectionParity` when layout or styling changes.
