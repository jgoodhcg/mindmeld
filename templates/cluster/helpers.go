package cluster

import (
	"fmt"
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
