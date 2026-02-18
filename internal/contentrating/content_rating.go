package contentrating

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Kids   int16 = 10
	Work   int16 = 20
	Adults int16 = 30
)

// Default returns the lobby content rating used when no explicit value is supplied.
func Default() int16 {
	return Work
}

func IsValid(id int16) bool {
	switch id {
	case Kids, Work, Adults:
		return true
	default:
		return false
	}
}

func ParseID(raw string) (int16, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return Default(), nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid content rating %q", raw)
	}

	id := int16(n)
	if !IsValid(id) {
		return 0, fmt.Errorf("unsupported content rating %d", n)
	}

	return id, nil
}

func Label(id int16) string {
	switch id {
	case Kids:
		return "Kids"
	case Work:
		return "Work"
	case Adults:
		return "Adults"
	default:
		return "Unknown"
	}
}

func Code(id int16) string {
	switch id {
	case Kids:
		return "kids"
	case Work:
		return "work"
	case Adults:
		return "adults"
	default:
		return "unknown"
	}
}
