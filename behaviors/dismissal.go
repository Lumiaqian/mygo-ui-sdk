package behaviors

import "github.com/egoist/mygo/ui"

// DismissReason identifies what triggered dismissal of an overlay or popover.
type DismissReason uint8

const (
	DismissNone DismissReason = iota
	DismissEscape
	DismissOutsidePress
)

// DismissOptions configures light-dismiss behavior for custom overlays.
type DismissOptions struct {
	// DisableEscape ignores Escape when true.
	DisableEscape bool
	// DisableOutsidePress ignores pointer presses outside surface when true.
	DisableOutsidePress bool
	// IgnoreEscapeWhileTyping skips Escape when a text input holds focus.
	IgnoreEscapeWhileTyping bool
}

// ShouldDismiss reports whether an open floating surface should close on this
// frame due to Escape or an outside pointer press.
func ShouldDismiss(c *ui.Context, surface ui.Element, opts DismissOptions) (bool, DismissReason) {
	if !opts.DisableEscape && (!opts.IgnoreEscapeWhileTyping || !c.FocusTakesText()) {
		if c.Shortcut(0, ui.KeyEscape) {
			return true, DismissEscape
		}
	}
	if !opts.DisableOutsidePress && surface.PressedOutside() {
		return true, DismissOutsidePress
	}
	return false, DismissNone
}
