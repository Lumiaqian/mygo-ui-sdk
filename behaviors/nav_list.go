package behaviors

import "github.com/egoist/mygo/ui"

// NavListOptions configures a headless tab-list navigation control such as a
// sidebar section switcher or top navigation bar.
type NavListOptions struct {
	Label       string
	Orientation ui.Orientation
	Wrap        bool
	Disabled    bool
}

// NavListParts holds the styled-by-caller RoleTabList container and per-item
// RoleTab constructor. Style Root before calling Root.Children(...).
type NavListParts struct {
	Root ui.Element
	sel  *int
	from int
	off  bool
}

// NavList creates a headless RoleTabList container with RoleTab items and
// arrow/Home/End keyboard switching (suppressed while c.FocusTakesText() is true).
// Default orientation is ui.Vertical; horizontal lists honor ui.RightToLeft.
func NavList(c *ui.Context, selected *int, count int, opts NavListOptions) NavListParts {
	orient := opts.Orientation
	if orient == 0 {
		orient = ui.Vertical
	}
	var fallback int
	if selected == nil {
		selected = &fallback
	}
	if count > 0 {
		*selected = clampInt(*selected, 0, count-1)
	}
	from := *selected

	var root ui.Element
	if orient == ui.Horizontal {
		root = ui.Row(c)
	} else {
		root = ui.Column(c)
	}
	root.Role(ui.RoleTabList).Orientation(orient)
	if opts.Label != "" {
		root.Label(opts.Label)
	}
	if opts.Disabled {
		root.Disabled(true)
	}

	if !opts.Disabled && count > 0 && !c.FocusTakesText() {
		stepBy := func(delta int) {
			next := *selected + delta
			if opts.Wrap {
				next = ((next % count) + count) % count
			} else {
				next = clampInt(next, 0, count-1)
			}
			*selected = next
		}
		if orient == ui.Horizontal {
			dir := 1
			if c.Direction() == ui.RightToLeft {
				dir = -1
			}
			switch {
			case root.Shortcut(0, ui.KeyRight):
				stepBy(dir)
			case root.Shortcut(0, ui.KeyLeft):
				stepBy(-dir)
			case root.Shortcut(0, ui.KeyHome):
				*selected = 0
			case root.Shortcut(0, ui.KeyEnd):
				*selected = count - 1
			}
		} else {
			switch {
			case root.Shortcut(0, ui.KeyDown):
				stepBy(1)
			case root.Shortcut(0, ui.KeyUp):
				stepBy(-1)
			case root.Shortcut(0, ui.KeyHome):
				*selected = 0
			case root.Shortcut(0, ui.KeyEnd):
				*selected = count - 1
			}
		}
	}

	return NavListParts{
		Root: root,
		sel:  selected,
		from: from,
		off:  opts.Disabled,
	}
}

// Item creates the index-th RoleTab button, marks AccessSelected when active,
// and updates the bound selection when clicked.
func (n NavListParts) Item(c *ui.Context, index int, label string) ui.Element {
	active := n.sel != nil && *n.sel == index
	el := ui.ButtonBase(c).
		Role(ui.RoleTab).
		Label(label).
		AccessSelected(active).
		Disabled(n.off)
	if !n.off && el.Clicked() && n.sel != nil {
		*n.sel = index
	}
	return el
}

// Selected returns the current selected index.
func (n NavListParts) Selected() int {
	if n.sel == nil {
		return 0
	}
	return *n.sel
}

// Changed reports whether the selection changed during this frame via keyboard or click.
func (n NavListParts) Changed() bool {
	return n.sel != nil && *n.sel != n.from
}
