package behaviors

import "github.com/egoist/mygo/ui"

// ToolbarOptions configures a headless toolbar with roving keyboard focus.
type ToolbarOptions struct {
	Label       string
	Orientation ui.Orientation
}

// ToolbarItemOptions configures optional accessibility and shortcut metadata on a toolbar item.
type ToolbarItemOptions struct {
	Tooltip       string
	Disabled      bool
	ShortcutMods  ui.Modifiers
	ShortcutKey   ui.Key
	HideFocusRing bool
}

// ToolbarParts holds the styled-by-caller RoleToolbar root and item constructors,
// following MyGo's immediate-mode *Parts convention so callers can style Root
// before invoking Root.Children(...).
type ToolbarParts struct {
	Root ui.Element
}

// Toolbar creates a headless RoleToolbar container with roving FocusGroup
// navigation (single Tab stop, arrow keys between items, Home/End to ends).
// Unlike ui.Toolbar, it imposes no visual styling or overflow chrome, making
// it suitable for custom floating or Liquid Glass toolbars.
func Toolbar(c *ui.Context, opts ToolbarOptions) ToolbarParts {
	orient := opts.Orientation
	if orient == 0 {
		orient = ui.Horizontal
	}
	var root ui.Element
	if orient == ui.Vertical {
		root = ui.Column(c)
	} else {
		root = ui.Row(c)
	}
	root.Role(ui.RoleToolbar).Orientation(orient).FocusGroup(orient)
	if opts.Label != "" {
		root.Label(opts.Label)
	}
	return ToolbarParts{Root: root}
}

// Button creates a headless focusable button inside the toolbar.
func (ToolbarParts) Button(c *ui.Context, label string, opts ToolbarItemOptions) ui.Element {
	el := ui.ButtonBase(c).Label(label).Disabled(opts.Disabled)
	return applyToolbarItemOptions(el, opts)
}

// Toggle creates a headless focusable toggle button (RoleToggleButton + Checked state).
func (ToolbarParts) Toggle(c *ui.Context, label string, active bool, opts ToolbarItemOptions) ui.Element {
	el := ui.ButtonBase(c).
		Role(ui.RoleToggleButton).
		Label(label).
		Checked(active).
		Disabled(opts.Disabled)
	return applyToolbarItemOptions(el, opts)
}

// RadioItem creates a headless focusable radio option (RoleRadio + Checked state),
// suitable for color swatches, stroke pickers, or mutually exclusive tool modes.
func (ToolbarParts) RadioItem(c *ui.Context, label string, active bool, opts ToolbarItemOptions) ui.Element {
	el := ui.ButtonBase(c).
		Role(ui.RoleRadio).
		Label(label).
		Checked(active).
		Disabled(opts.Disabled)
	return applyToolbarItemOptions(el, opts)
}

func applyToolbarItemOptions(el ui.Element, opts ToolbarItemOptions) ui.Element {
	if opts.Tooltip != "" {
		el.Tooltip(opts.Tooltip)
	}
	if opts.ShortcutKey != 0 {
		el.KeyShortcuts(opts.ShortcutMods, opts.ShortcutKey)
	}
	if opts.HideFocusRing {
		el.FocusRing(false)
	}
	return el
}
