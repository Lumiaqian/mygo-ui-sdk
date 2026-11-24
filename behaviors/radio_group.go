package behaviors

import "github.com/egoist/mygo/ui"

// RadioGroupOptions configures a headless radio group with roving FocusGroup
// navigation (suitable for color swatches, template cards, or custom choice chips).
type RadioGroupOptions struct {
	Label       string
	Orientation ui.Orientation
	Disabled    bool
}

// RadioItemOptions configures optional metadata on an item inside RadioGroup.
type RadioItemOptions struct {
	Tooltip       string
	Disabled      bool
	HideFocusRing bool
}

// RadioGroupParts holds the styled-by-caller RoleRadioGroup container and
// RoleRadio item builder. Style Root before calling Root.Children(...).
type RadioGroupParts struct {
	Root ui.Element
	off  bool
}

// RadioGroup creates a headless RoleRadioGroup container with FocusGroup roving
// focus (single Tab stop landing on the checked item, arrow keys and Home/End
// between choices). Default orientation is ui.Horizontal.
func RadioGroup(c *ui.Context, opts RadioGroupOptions) RadioGroupParts {
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
	root.Role(ui.RoleRadioGroup).Orientation(orient).FocusGroup(orient)
	if opts.Label != "" {
		root.Label(opts.Label)
	}
	if opts.Disabled {
		root.Disabled(true)
	}
	return RadioGroupParts{
		Root: root,
		off:  opts.Disabled,
	}
}

// Item creates a headless focusable RoleRadio option with Checked state.
func (g RadioGroupParts) Item(c *ui.Context, label string, active bool, opts RadioItemOptions) ui.Element {
	el := ui.ButtonBase(c).
		Role(ui.RoleRadio).
		Label(label).
		Checked(active).
		Disabled(g.off || opts.Disabled)
	if opts.Tooltip != "" {
		el.Tooltip(opts.Tooltip)
	}
	if opts.HideFocusRing {
		el.FocusRing(false)
	}
	return el
}
