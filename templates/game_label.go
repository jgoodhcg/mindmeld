package templates

import "strings"

func gameLabel(slug string) string {
	switch strings.ToLower(strings.TrimSpace(slug)) {
	case "trivia":
		return "Trivia"
	case "cluster":
		return "Cluster"
	default:
		if slug == "" {
			return "Game"
		}
		return strings.ToUpper(slug[:1]) + strings.ToLower(slug[1:])
	}
}
