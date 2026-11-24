package behaviors

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// SegmentedTrackOptions configures a single-track draggable segmented control
// (for example, a Liquid Glass sliding lens segmented control).
type SegmentedTrackOptions struct {
	Label         string
	Labels        []string
	SegmentWidth  float32
	Padding       float32
	Disabled      bool
	DragThreshold float32
	FrameInterval time.Duration
}

// SegmentedTrackResult returns the focusable track element and current lens state.
type SegmentedTrackResult struct {
	Track    ui.Element
	Changed  bool
	Dragging bool
	Target   float32
	Selected int
}

// SegmentedTrack creates a continuous draggable segmented track with full keyboard
// navigation (Left/Right honoring RTL, Home, End, and Enter/Space cycle),
// drag-to-snap, click-to-segment, and accessible Label + Value semantics.
// All visual rendering (track background, sliding lens, labels) stays with the caller.
func SegmentedTrack(c *ui.Context, selected *int, opts SegmentedTrackOptions) SegmentedTrackResult {
	n := len(opts.Labels)
	if n == 0 {
		return SegmentedTrackResult{Track: ui.Row(c)}
	}
	if *selected < 0 {
		*selected = 0
	} else if *selected >= n {
		*selected = n - 1
	}

	from := *selected
	hi := float32(n-1) * opts.SegmentWidth
	restPos := float32(*selected) * opts.SegmentWidth

	track := ui.ButtonBase(c).
		Label(opts.Label).
		Value(opts.Labels[*selected]).
		Orientation(ui.Horizontal).
		Disabled(opts.Disabled)

	drag := TrackDrag(c, track, "behaviors-segmented-drag", restPos, opts.Disabled, TrackDragOptions{
		Min:           0,
		Max:           hi,
		OriginOffset:  opts.Padding,
		Threshold:     opts.DragThreshold,
		FrameInterval: opts.FrameInterval,
	})

	if !opts.Disabled {
		switch {
		case drag.Released && drag.WasDrag && opts.SegmentWidth > 0:
			idx := int(drag.Pos/opts.SegmentWidth + 0.5)
			*selected = clampInt(idx, 0, n-1)
		case drag.Released && !drag.WasDrag && opts.SegmentWidth > 0:
			idx := int(drag.Origin / opts.SegmentWidth)
			*selected = clampInt(idx, 0, n-1)
		case !drag.Active && track.Clicked():
			*selected = (*selected + 1) % n
		}

		if track.Focused() {
			step := 1
			if c.Direction() == ui.RightToLeft {
				step = -1
			}
			switch {
			case track.Shortcut(0, ui.KeyRight):
				*selected = clampInt(*selected+step, 0, n-1)
			case track.Shortcut(0, ui.KeyLeft):
				*selected = clampInt(*selected-step, 0, n-1)
			case track.Shortcut(0, ui.KeyHome):
				*selected = 0
			case track.Shortcut(0, ui.KeyEnd):
				*selected = n - 1
			}
		}
	}

	track.Value(opts.Labels[*selected])

	target := float32(*selected) * opts.SegmentWidth
	if drag.Active {
		target = drag.Pos
	}

	return SegmentedTrackResult{
		Track:    track,
		Changed:  *selected != from,
		Dragging: drag.Active,
		Target:   target,
		Selected: *selected,
	}
}

// SegmentedGroupOptions configures a discrete segmented control built on ui.SegmentedBase.
type SegmentedGroupOptions struct {
	Label    string
	Labels   []string
	Disabled bool
}

// SegmentedGroupParts wraps ui.SegmentedParts with a named track and item helper.
type SegmentedGroupParts struct {
	Track  ui.Element
	parts  ui.SegmentedParts
	labels []string
	sel    *int
	from   int
	off    bool
}

// SegmentedGroup creates a headless discrete segmented control on top of
// ui.SegmentedBase so custom-styled segmented controls inherit RoleRadioGroup,
// RoleRadio segments, FocusGroup arrow navigation, and accessible labels.
func SegmentedGroup(c *ui.Context, selected *int, opts SegmentedGroupOptions) SegmentedGroupParts {
	n := len(opts.Labels)
	if n > 0 {
		*selected = clampInt(*selected, 0, n-1)
	}
	from := *selected
	parts := ui.SegmentedBase(c, selected, n)
	if opts.Label != "" {
		parts.Track.Label(opts.Label)
	}
	if opts.Disabled {
		parts.Track.Disabled(true)
	}
	return SegmentedGroupParts{
		Track:  parts.Track,
		parts:  parts,
		labels: opts.Labels,
		sel:    selected,
		from:   from,
		off:    opts.Disabled,
	}
}

// Segment returns the i-th accessible segment element with its label pre-applied.
func (g SegmentedGroupParts) Segment(i int) ui.Element {
	el := g.parts.Segment(i)
	if i >= 0 && i < len(g.labels) {
		el.Label(g.labels[i])
	}
	if g.off {
		el.Disabled(true)
	}
	return el
}

// Changed reports whether the selection changed during this frame.
func (g SegmentedGroupParts) Changed() bool {
	return g.sel != nil && *g.sel != g.from
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
