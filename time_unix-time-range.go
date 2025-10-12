// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"context"
	stdtime "time"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

type UnixTimeRanges []UnixTimeRange

// Max returns the maximum UnixTimeRange that encompasses all ranges in the list.
// It finds the earliest From time and the latest Until time across all ranges.
// Returns nil if the list is empty.
func (ranges UnixTimeRanges) Max() *UnixTimeRange {
	if len(ranges) == 0 {
		return nil
	}

	maxRange := ranges[0]
	for _, r := range ranges[1:] {
		if r.From.Before(maxRange.From) {
			maxRange.From = r.From
		}
		if r.Until.After(maxRange.Until) {
			maxRange.Until = r.Until
		}
	}

	return maxRange.Ptr()
}

// Min returns the minimum UnixTimeRange that is contained within all ranges in the list.
// It finds the latest From time and the earliest Until time across all ranges.
// Returns nil if the list is empty or if there is no overlap between ranges.
func (ranges UnixTimeRanges) Min() *UnixTimeRange {
	if len(ranges) == 0 {
		return nil
	}

	minRange := ranges[0]
	for _, r := range ranges[1:] {
		if r.From.After(minRange.From) {
			minRange.From = r.From
		}
		if r.Until.Before(minRange.Until) {
			minRange.Until = r.Until
		}
	}

	// Check if the resulting range is valid (From <= Until)
	if minRange.From.After(minRange.Until) {
		return nil
	}

	return minRange.Ptr()
}

// UnixTimeRangeFromTime creates a UnixTimeRange from two time.Time values.
// It converts the from and until times to UnixTime types and returns a UnixTimeRange.
func UnixTimeRangeFromTime(from, until stdtime.Time) UnixTimeRange {
	return UnixTimeRange{
		From:  UnixTime(from),
		Until: UnixTime(until),
	}
}

type UnixTimeRange struct {
	From  UnixTime `json:"from,omitempty"`
	Until UnixTime `json:"until,omitempty"`
}

func (r UnixTimeRange) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("from", r.From),
		validation.Name("until", r.Until),
		validation.Name("range", validation.HasValidationFunc(func(ctx context.Context) error {
			if r.From.After(r.Until) {
				return errors.Wrapf(
					ctx,
					validation.Error,
					"from must be less than or equal to until",
				)
			}
			return nil
		})),
	}.Validate(ctx)
}

func (r UnixTimeRange) Ptr() *UnixTimeRange {
	return &r
}

// DayUnixTimeRange creates a UnixTimeRange covering the entire day containing the given unix time.
// The range spans from 00:00:00.000000000 to 23:59:59.999999999 of that day.
func DayUnixTimeRange(ut UnixTime) UnixTimeRange {
	t := ut.Time()
	return UnixTimeRange{From: UnixTime(BeginningOfDay(t)), Until: UnixTime(EndOfDay(t))}
}

// WeekUnixTimeRange creates a UnixTimeRange covering the entire week containing the given unix time.
// The range spans from Monday 00:00:00.000000000 to Sunday 23:59:59.999999999 of that week.
// Uses ISO 8601 standard where Monday is the first day of the week.
func WeekUnixTimeRange(ut UnixTime) UnixTimeRange {
	t := ut.Time()
	return UnixTimeRange{From: UnixTime(BeginningOfWeek(t)), Until: UnixTime(EndOfWeek(t))}
}

// MonthUnixTimeRange creates a UnixTimeRange covering the entire month containing the given unix time.
// The range spans from the 1st day 00:00:00.000000000 to the last day 23:59:59.999999999 of that month.
func MonthUnixTimeRange(ut UnixTime) UnixTimeRange {
	t := ut.Time()
	return UnixTimeRange{From: UnixTime(BeginningOfMonth(t)), Until: UnixTime(EndOfMonth(t))}
}

// QuarterUnixTimeRange creates a UnixTimeRange covering the entire quarter containing the given unix time.
// Quarters are defined as: Q1=Jan-Mar, Q2=Apr-Jun, Q3=Jul-Sep, Q4=Oct-Dec.
// The range spans from the 1st day of the quarter 00:00:00.000000000 to the last day 23:59:59.999999999 of that quarter.
func QuarterUnixTimeRange(ut UnixTime) UnixTimeRange {
	t := ut.Time()
	return UnixTimeRange{From: UnixTime(BeginningOfQuarter(t)), Until: UnixTime(EndOfQuarter(t))}
}

// YearUnixTimeRange creates a UnixTimeRange covering the entire year containing the given unix time.
// The range spans from January 1st 00:00:00.000000000 to December 31st 23:59:59.999999999 of that year.
func YearUnixTimeRange(ut UnixTime) UnixTimeRange {
	t := ut.Time()
	return UnixTimeRange{From: UnixTime(BeginningOfYear(t)), Until: UnixTime(EndOfYear(t))}
}

// TimeRange converts a UnixTimeRange to a TimeRange.
func (r UnixTimeRange) TimeRange() TimeRange {
	return TimeRange{From: r.From.Time(), Until: r.Until.Time()}
}
