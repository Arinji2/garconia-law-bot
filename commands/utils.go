package commands_utils

import (
	"fmt"
	"strconv"
)

func FormatDescription(description string) string {
	if len(description) > 90 {
		return description[:90] + "..."
	}
	return description
}

func OrdinalRepresentation(s string) string {
	n, err := strconv.Atoi(s)
	if err != nil {
		return s
	}
	if n <= 0 {
		return fmt.Sprintf("%d", n) // Return just the number for non-positive values
	}

	// Special cases for numbers ending in 11, 12, or 13
	if n%100 >= 11 && n%100 <= 13 {
		return fmt.Sprintf("%dth", n)
	}

	suffix := "th"
	switch n % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	}

	return fmt.Sprintf("%d%s", n, suffix)
}
