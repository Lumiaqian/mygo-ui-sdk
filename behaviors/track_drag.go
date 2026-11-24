package behaviors

import (
	"time"

	"github.com/egoist/mygo/ui"
)

const (
	defaultDragThreshold = float32(3)
	defaultDragFrame     = 16 * time.Millisecond
)

// TrackDragOptions configures 1D pointer drag tracking with a movement threshold.
type TrackDragOptions struct {
	Min           float32
	Max           float32
	OriginOffset  float32
	Threshold     float32
	FrameInterval time.Duration
}

// TrackDragResult reports the current 1D drag state for the frame.
type TrackDragResult struct {
	Active   bool
	Pos      float32
	Moved    float32
	StartPos float32
	Origin   float32
	Released bool
	WasDrag  bool
}

type trackDragState struct {
	active   bool
	pos      float32
	moved    float32
	startPos float32
	origin   float32
}

// TrackDrag manages a 1D pointer drag on el using ui.Local state:
// it tracks pointer delta clamped to [opts.Min, opts.Max], distinguishes taps
// from drags using opts.Threshold, and schedules frames while active so release
// animations settle reliably.
func TrackDrag(c *ui.Context, el ui.Element, key any, restPos float32, disabled bool, opts TrackDragOptions) TrackDragResult {
	threshold := opts.Threshold
	if threshold <= 0 {
		threshold = defaultDragThreshold
	}
	interval := opts.FrameInterval
	if interval <= 0 {
		interval = defaultDragFrame
	}

	d := ui.Local(el, key, func() trackDragState { return trackDragState{} })
	px, _, _ := el.PointerPosition()

	if dx, _, pressed := el.Dragged(); pressed && !disabled {
		if !d.active {
			*d = trackDragState{
				active:   true,
				pos:      clampf(restPos, opts.Min, opts.Max),
				startPos: restPos,
				origin:   px - opts.OriginOffset,
			}
		}
		d.pos = clampf(d.pos+dx, opts.Min, opts.Max)
		d.moved += absf(dx)
		c.After(interval)
		return TrackDragResult{
			Active:   true,
			Pos:      d.pos,
			Moved:    d.moved,
			StartPos: d.startPos,
			Origin:   d.origin,
		}
	}

	if d.active {
		res := TrackDragResult{
			Active:   false,
			Pos:      d.pos,
			Moved:    d.moved,
			StartPos: d.startPos,
			Origin:   d.origin,
			Released: true,
			WasDrag:  d.moved > threshold,
		}
		d.active = false
		return res
	}

	return TrackDragResult{Pos: clampf(restPos, opts.Min, opts.Max)}
}

// SwitchTrackOptions configures a draggable headless switch built on ui.SwitchBase.
type SwitchTrackOptions struct {
	Label         string
	Min           float32
	Max           float32
	Disabled      bool
	DragThreshold float32
	FrameInterval time.Duration
}

// SwitchTrackResult holds the underlying ui.SwitchBase element and knob drag state.
type SwitchTrackResult struct {
	Switch   ui.Element
	Changed  bool
	Active   bool
	Target   float32
	Progress float32
}

type switchDragStart struct {
	on bool
}

// SwitchTrack wraps ui.SwitchBase with 1D knob dragging and snap-to-side on release.
// Keyboard toggle (Space/Enter) and RoleSwitch semantics come from ui.SwitchBase.
func SwitchTrack(c *ui.Context, on *bool, opts SwitchTrackOptions) SwitchTrackResult {
	side := func(v bool) float32 {
		if v {
			return opts.Max
		}
		return opts.Min
	}

	sw := ui.SwitchBase(c, on).Label(opts.Label).Disabled(opts.Disabled)
	changed := sw.Changed()
	start := ui.Local(sw, "behaviors-switch-start", func() switchDragStart { return switchDragStart{} })

	rest := side(*on)
	drag := TrackDrag(c, sw, "behaviors-switch-drag", rest, opts.Disabled, TrackDragOptions{
		Min:           opts.Min,
		Max:           opts.Max,
		Threshold:     opts.DragThreshold,
		FrameInterval: opts.FrameInterval,
	})

	if drag.Active && drag.Moved == 0 {
		start.on = *on
	}

	target := rest
	if drag.Active {
		target = drag.Pos
	} else if drag.Released && drag.WasDrag {
		mid := (opts.Min + opts.Max) * 0.5
		*on = drag.Pos > mid
		changed = *on != start.on
		target = side(*on)
	}

	var progress float32
	if span := opts.Max - opts.Min; span > 0 {
		progress = clampf((target-opts.Min)/span, 0, 1)
	}
	return SwitchTrackResult{
		Switch:   sw,
		Changed:  changed,
		Active:   drag.Active,
		Target:   target,
		Progress: progress,
	}
}

func clampf(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func absf(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}
