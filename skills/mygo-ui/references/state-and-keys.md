# Immediate-Mode State, Keys, and Focus

## 1. `Context.Key` vs `Element.Key`

- Stateful bases (`SelectBase`, `SegmentedBase`, `TabsBase`, `TreeBase`, `MenuBase`) and any widgets inside loops MUST receive a stable key via `c.Key(...)` **before** construction.
- Stable keys preserve `ui.Local`, `AnimateTo`, focus, scroll offset, and automation IDs across frames when sibling order changes.
- Key loop items by semantic ID (e.g., `c.Key("color-"+col.ID)`), not fragile slice index when items can reorder or hide.
- **Always key containers that follow conditional siblings** (such as floating toolbars, HUDs, or preview layers inside an overlay root). Because MyGo derives child IDs from `keyedID(parentID, key)`, an unkeyed ancestor shifts every descendant's ID whenever an earlier conditional sibling mounts or unmounts—dropping focus, breaking hover, and swallowing clicks even when children themselves are keyed.

## 2. `ui.Local` and cleanup

- Store transient interaction state (drag state, hover timers, local animation bookkeeping) with `ui.Local(el, "my-key", initFn)`.
- If the stored value implements `io.Closer`, MyGo calls `Close()` once on the main thread when the element leaves the tree. Keep `Close()` non-blocking (cancel a context, do not wait for goroutines).

## 3. Ephemeral values

- `ui.Element` and `*Base` part structs expire at the end of the build pass. Never store them on structs or pass them to goroutines.
- Capture `c.Services()` for background invalidation (`Services.Invalidate()`) and `ui.Handle` / `FocusBind` for persistent focus control.

## 4. Drag release animation frames

- Releasing a drag does not emit a click event, so an element animating back to its resting position can stall unless the drag frame schedules a follow-up frame via `c.After(16 * time.Millisecond)`. `behaviors.TrackDrag` handles this automatically.
