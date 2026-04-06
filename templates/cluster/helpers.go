package cluster

import (
	"fmt"
	"math"
	"strings"
)

func plotStyle(x, y float64, radiusPx int, withAnimation bool) string {
	clampedX := clampUnit(x)
	clampedY := clampUnit(y)
	left := clampedX * 100
	top := (1 - clampedY) * 100
	style := fmt.Sprintf("left: calc(%.2f%% - %dpx); top: calc(%.2f%% - %dpx);", left, radiusPx, top, radiusPx)
	if withAnimation {
		style += " animation: cluster-dot-enter 260ms ease-out both;"
	}
	return style
}

func plotStyleAnimated(x, y float64, radiusPx int, delayMs int) string {
	style := plotStyle(x, y, radiusPx, true)
	if delayMs > 0 {
		style += fmt.Sprintf(" animation-delay: %dms;", delayMs)
	}
	return style
}

func dotClass(dot DotView) string {
	className := "absolute z-10 h-4 w-4 rounded-full border border-base bg-text"
	if dot.IsCurrentPlayer {
		className += " ring-2 ring-cyan ring-offset-2 ring-offset-base"
	}
	return className
}

func dotLabelClass(dot DotView) string {
	className := "absolute z-0 whitespace-nowrap rounded bg-base/75 px-1.5 py-0.5 text-[10px] leading-none text-text-muted"
	if dot.IsOutlier {
		className += " border border-amber/40 text-amber"
	}
	return className
}

func dotLabelStyle(dot DotView) string {
	left := clampUnit(dot.X) * 100
	top := (1 - clampUnit(dot.Y)) * 100

	xTransform := "-50%"
	textAlign := "center"
	switch {
	case dot.X <= 0.18:
		xTransform = "0"
		textAlign = "left"
	case dot.X >= 0.82:
		xTransform = "-100%"
		textAlign = "right"
	}

	labelTop := fmt.Sprintf("calc(%.2f%% - 14px)", top)
	if dot.Y >= 0.82 {
		labelTop = fmt.Sprintf("calc(%.2f%% + 12px)", top)
	}

	return fmt.Sprintf("left: %.2f%%; top: %s; transform: translateX(%s); text-align: %s;", left, labelTop, xTransform, textAlign)
}

func displayCoord(v float64) string {
	dv := (clampUnit(v) - 0.5) * 2
	if math.Abs(dv) < 0.005 {
		dv = 0
	}
	return fmt.Sprintf("%.2f", dv)
}

func submissionProgress(submittedCount int, expectedCount int) string {
	return fmt.Sprintf("%d / %d submitted", submittedCount, expectedCount)
}

func formatAveragePoints(points float64) string {
	return fmt.Sprintf("%.1f", points)
}

func formatDistance(distance float64) string {
	return fmt.Sprintf("%.2f", distance)
}

func centerStyle(x, y float64) string {
	clampedX := clampUnit(x)
	clampedY := clampUnit(y)
	left := clampedX * 100
	top := (1 - clampedY) * 100
	return fmt.Sprintf("left: %.2f%%; top: %.2f%%;", left, top)
}

func outlierSummary(outliers []string) string {
	switch len(outliers) {
	case 0:
		return ""
	case 1:
		return outliers[0] + " landed furthest from the group center"
	default:
		return strings.Join(outliers, ", ") + " landed furthest from the group center"
	}
}

func spreadRingStyle(x, y float64) string {
	return centerStyle(x, y) + "width: 25%; height: 25%; transform: translate(-50%, -50%);"
}

func clampUnit(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
