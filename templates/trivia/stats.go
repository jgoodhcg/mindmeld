package trivia

import "math"

// CalculatePercentage computes the percentage of count in total.
func CalculatePercentage(count int, total int) int {
	if total == 0 {
		return 0
	}
	return int(math.Round((float64(count) / float64(total)) * 100))
}
