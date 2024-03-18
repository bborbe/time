package time

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/bborbe/errors"
)

const (
	Nanosecond  time.Duration = 1
	Microsecond               = 1000 * Nanosecond
	Millisecond               = 1000 * Microsecond
	Second                    = 1000 * Millisecond
	Minute                    = 60 * Second
	Hour                      = 60 * Minute
	Day                       = 24 * Hour
	Week                      = 7 * Day
)

// UnitMap contains units to duration mapping
var UnitMap = map[string]time.Duration{
	"ns": Nanosecond,
	"us": Microsecond,
	"µs": Microsecond, // U+00B5 = micro symbol
	"μs": Microsecond, // U+03BC = Greek letter mu
	"ms": Millisecond,
	"s":  Second,
	"m":  Minute,
	"h":  Hour,
	"d":  Day,
	"w":  Week,
}

var durationRegexp = regexp.MustCompile(`(\d*\.?\d+)([a-z]+)`)

func ParseDuration(ctx context.Context, value string) (*time.Duration, error) {
	var isNegative bool
	if len(value) > 0 && value[0] == '-' {
		isNegative = true
		value = value[1:]
	}
	var result time.Duration
	for _, match := range durationRegexp.FindAllStringSubmatch(value, -1) {
		if len(match) != 3 {
			return nil, errors.Errorf(ctx, "invalid length of match")
		}
		value, err := parseAsDuration(ctx, match[1], match[2])
		if err != nil {
			return nil, errors.Wrapf(ctx, err, "parse failed")
		}
		result += value
	}
	if isNegative {
		result = result * -1
	}
	return &result, nil
}

func parseAsDuration(ctx context.Context, value string, unit string) (time.Duration, error) {
	factor, ok := UnitMap[unit]
	if !ok {
		return 0, errors.Errorf(ctx, "unkown unit '%s'", unit)
	}
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, errors.Wrapf(ctx, err, "parse failed")
	}
	return time.Duration(i * float64(factor)), nil
}
