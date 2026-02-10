package cluster

import "testing"

func TestCalculateCentroid(t *testing.T) {
	points := []Point{
		{X: 0.2, Y: 0.8},
		{X: 0.8, Y: 0.2},
		{X: 0.5, Y: 0.5},
	}

	x, y, ok := CalculateCentroid(points)
	if !ok {
		t.Fatal("expected centroid calculation to succeed")
	}

	if x != 0.5 {
		t.Fatalf("expected centroid x=0.5, got %v", x)
	}
	if y != 0.5 {
		t.Fatalf("expected centroid y=0.5, got %v", y)
	}
}

func TestCalculateRoundPoints(t *testing.T) {
	centroidX := 0.5
	centroidY := 0.5

	if got := CalculateRoundPoints(0.5, 0.5, centroidX, centroidY); got != 100 {
		t.Fatalf("expected perfect match to score 100, got %d", got)
	}

	if got := CalculateRoundPoints(1.0, 1.0, 0.0, 0.0); got != 0 {
		t.Fatalf("expected max distance to score 0, got %d", got)
	}

	if got := CalculateRoundPoints(0.5, 1.0, 0.5, 0.5); got != 65 {
		t.Fatalf("expected midpoint offset score 65, got %d", got)
	}
}

func TestSelectNextUnusedPairNoRepeat(t *testing.T) {
	ordered := []string{"pair-a", "pair-b", "pair-c"}
	used := map[string]bool{"pair-a": true}

	next, ok := SelectNextUnusedPair(ordered, used)
	if !ok {
		t.Fatal("expected next unused pair")
	}
	if next != "pair-b" {
		t.Fatalf("expected pair-b, got %s", next)
	}

	used[next] = true
	next, ok = SelectNextUnusedPair(ordered, used)
	if !ok || next != "pair-c" {
		t.Fatalf("expected pair-c after marking pair-b used, got %s (ok=%v)", next, ok)
	}
}

func TestIsPairPoolExhausted(t *testing.T) {
	ordered := []string{"pair-a", "pair-b"}

	if IsPairPoolExhausted(ordered, map[string]bool{"pair-a": true}) {
		t.Fatal("expected pool to still have one pair left")
	}

	if !IsPairPoolExhausted(ordered, map[string]bool{"pair-a": true, "pair-b": true}) {
		t.Fatal("expected pool to be exhausted")
	}
}
