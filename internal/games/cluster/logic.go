package cluster

import "math"

const (
	minPlayersToStart = 2
	maxDistance       = math.Sqrt2
)

// Point represents a normalized point in the unit square.
type Point struct {
	X float64
	Y float64
}

// CalculateCentroid returns the average x/y for a set of points.
func CalculateCentroid(points []Point) (float64, float64, bool) {
	if len(points) == 0 {
		return 0, 0, false
	}

	sumX := 0.0
	sumY := 0.0
	for _, p := range points {
		sumX += clampUnit(p.X)
		sumY += clampUnit(p.Y)
	}

	count := float64(len(points))
	return sumX / count, sumY / count, true
}

// CalculateRoundPoints applies the centroid-distance scoring formula.
func CalculateRoundPoints(x, y, centroidX, centroidY float64) int {
	distance := CalculateDistance(x, y, centroidX, centroidY)
	normalized := math.Min(distance/maxDistance, 1)
	return int(math.Round((1 - normalized) * 100))
}

// CalculateDistance returns Euclidean distance between two normalized points.
func CalculateDistance(x1, y1, x2, y2 float64) float64 {
	dx := clampUnit(x1) - clampUnit(x2)
	dy := clampUnit(y1) - clampUnit(y2)
	return math.Sqrt((dx * dx) + (dy * dy))
}

// SelectNextUnusedPair returns the next available pair id from a deterministic ordered list.
func SelectNextUnusedPair(orderedPairIDs []string, used map[string]bool) (string, bool) {
	for _, id := range orderedPairIDs {
		if !used[id] {
			return id, true
		}
	}
	return "", false
}

// IsPairPoolExhausted reports whether all ordered ids are already marked used.
func IsPairPoolExhausted(orderedPairIDs []string, used map[string]bool) bool {
	_, ok := SelectNextUnusedPair(orderedPairIDs, used)
	return !ok
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
