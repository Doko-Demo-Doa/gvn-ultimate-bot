package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDuration parses human-readable duration strings like "30m", "2h", "7d", "1w".
// It falls back to time.ParseDuration for standard Go duration strings.
func ParseDuration(input string) (time.Duration, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Try standard Go duration first (handles ns, us, ms, s, m, h)
	if d, err := time.ParseDuration(input); err == nil {
		return d, nil
	}

	if len(input) < 2 {
		return 0, fmt.Errorf("invalid duration format: %s", input)
	}

	unit := input[len(input)-1]
	valueStr := input[:len(input)-1]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %s", input)
	}

	switch unit {
	case 'd':
		return time.Duration(value * float64(24*time.Hour)), nil
	case 'w':
		return time.Duration(value * float64(7*24*time.Hour)), nil
	default:
		return 0, fmt.Errorf("unsupported duration unit %q in %q (supported: m, h, d, w)", string(unit), input)
	}
}
