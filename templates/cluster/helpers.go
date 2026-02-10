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

func dotClass(dot DotView) string {
	// Default all other players to gray; highlight current player in cyan.
	className := "absolute h-4 w-4 rounded-full border border-base bg-text-muted"
	if dot.IsCurrentPlayer {
		className = "absolute h-4 w-4 rounded-full border border-base bg-cyan"
	}

	// Winner outline matches amber scoreboard highlight.
	if dot.IsWinner {
		className += " ring-2 ring-amber ring-offset-1 ring-offset-base"
	}

	return className
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

func averageTotalPoints(standings []StandingView) float64 {
	if len(standings) == 0 {
		return 0
	}

	sum := 0
	for _, s := range standings {
		sum += s.TotalPoints
	}

	return float64(sum) / float64(len(standings))
}

func formatAveragePoints(points float64) string {
	return fmt.Sprintf("%.1f", points)
}

func centerStyle(x, y float64) string {
	clampedX := clampUnit(x)
	clampedY := clampUnit(y)
	left := clampedX * 100
	top := (1 - clampedY) * 100
	return fmt.Sprintf("left: %.2f%%; top: %.2f%%;", left, top)
}

func winnerSummary(winners []string) string {
	switch len(winners) {
	case 0:
		return ""
	case 1:
		return winners[0] + " wins this round"
	default:
		return strings.Join(winners, ", ") + " tie this round"
	}
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
