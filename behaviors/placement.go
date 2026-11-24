package behaviors

import (
	"image"

	"github.com/egoist/mygo/ui"
)

// PlacementSide selects which side of an anchor a floating surface prefers.
type PlacementSide uint8

const (
	SideBottom PlacementSide = iota
	SideTop
	SideRight
	SideLeft
)

// PlacementAlign aligns a floating surface along the anchor's cross axis.
type PlacementAlign uint8

const (
	AlignStart PlacementAlign = iota
	AlignCenter
	AlignEnd
)

// PlacementOptions configures floating surface positioning relative to an anchor.
type PlacementOptions struct {
	PreferredSide PlacementSide
	Align         PlacementAlign
	Gap           int
	Margin        int
	Direction     ui.Direction
	// InsideFallback places the surface inside the anchor edge when neither the
	// preferred nor opposite side fits inside the viewport.
	InsideFallback bool
}

// PlacementResult is the resolved viewport-clamped position and chosen side.
type PlacementResult struct {
	Point       image.Point
	Rect        image.Rectangle
	Side        PlacementSide
	Flipped     bool
	AnchoredEnd bool
}

// PlaceFloating computes a viewport-safe position for a floating surface of
// size next to anchor inside viewport. It flips to the opposite side when the
// preferred side overflows, optionally falls back inside the anchor, clamps to
// viewport margins, and mirrors Start/End alignment under ui.RightToLeft.
func PlaceFloating(viewport, anchor image.Rectangle, size image.Point, opts PlacementOptions) PlacementResult {
	gap := opts.Gap
	margin := opts.Margin
	if margin <= 0 {
		margin = gap
	}

	align := opts.Align
	if opts.Direction == ui.RightToLeft {
		switch align {
		case AlignStart:
			align = AlignEnd
		case AlignEnd:
			align = AlignStart
		}
	}

	side := opts.PreferredSide
	flipped := false
	var x, y int

	switch side {
	case SideBottom, SideTop:
		x = alignSpan(anchor.Min.X, anchor.Max.X, size.X, align)
		x = clampAxis(x, viewport.Min.X+margin, viewport.Max.X-size.X-margin)

		y, side, flipped = placeVertical(viewport, anchor, size.Y, gap, margin, side, opts.InsideFallback)
		y = clampAxis(y, viewport.Min.Y+margin, viewport.Max.Y-size.Y-margin)

	case SideRight, SideLeft:
		y = alignSpan(anchor.Min.Y, anchor.Max.Y, size.Y, align)
		y = clampAxis(y, viewport.Min.Y+margin, viewport.Max.Y-size.Y-margin)

		x, side, flipped = placeHorizontal(viewport, anchor, size.X, gap, margin, side, opts.InsideFallback)
		x = clampAxis(x, viewport.Min.X+margin, viewport.Max.X-size.X-margin)
	}

	pt := image.Pt(x, y)
	return PlacementResult{
		Point:       pt,
		Rect:        image.Rectangle{Min: pt, Max: pt.Add(size)},
		Side:        side,
		Flipped:     flipped,
		AnchoredEnd: align == AlignEnd && anchor.Max.X >= min(viewport.Max.X-margin, size.X+margin*2),
	}
}

func placeVertical(vp, anchor image.Rectangle, h, gap, margin int, pref PlacementSide, insideFallback bool) (int, PlacementSide, bool) {
	below := anchor.Max.Y + gap
	above := anchor.Min.Y - h - gap
	fitsBelow := below+h+margin <= vp.Max.Y
	fitsAbove := above >= vp.Min.Y+margin

	if pref == SideBottom {
		if fitsBelow {
			return below, SideBottom, false
		}
		if fitsAbove {
			return above, SideTop, true
		}
		if insideFallback {
			return max(vp.Min.Y+margin, anchor.Max.Y-h-gap), SideBottom, false
		}
		return below, SideBottom, false
	}

	if fitsAbove {
		return above, SideTop, false
	}
	if fitsBelow {
		return below, SideBottom, true
	}
	if insideFallback {
		return min(vp.Max.Y-h-margin, anchor.Min.Y+gap), SideTop, false
	}
	return above, SideTop, false
}

func placeHorizontal(vp, anchor image.Rectangle, w, gap, margin int, pref PlacementSide, insideFallback bool) (int, PlacementSide, bool) {
	right := anchor.Max.X + gap
	left := anchor.Min.X - w - gap
	fitsRight := right+w+margin <= vp.Max.X
	fitsLeft := left >= vp.Min.X+margin

	if pref == SideRight {
		if fitsRight {
			return right, SideRight, false
		}
		if fitsLeft {
			return left, SideLeft, true
		}
		if insideFallback {
			return max(vp.Min.X+margin, anchor.Max.X-w-gap), SideRight, false
		}
		return right, SideRight, false
	}

	if fitsLeft {
		return left, SideLeft, false
	}
	if fitsRight {
		return right, SideRight, true
	}
	if insideFallback {
		return min(vp.Max.X-w-margin, anchor.Min.X+gap), SideLeft, false
	}
	return left, SideLeft, false
}

func alignSpan(aMin, aMax, size int, align PlacementAlign) int {
	switch align {
	case AlignCenter:
		return (aMin+aMax)/2 - size/2
	case AlignEnd:
		return aMax - size
	default:
		return aMin
	}
}

func clampAxis(v, lo, hi int) int {
	if hi < lo {
		return lo
	}
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
