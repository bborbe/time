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

type DateRanges []DateRange

// Max returns the maximum DateRange that encompasses all ranges in the list.
// It finds the earliest From date and the latest Until date across all ranges.
// Returns nil if the list is empty.
func (ranges DateRanges) Max() *DateRange {
	if len(ranges) == 0 {
		return nil
	}

	maxRange := ranges[0]
	for _, r := range ranges[1:] {
		if r.From.Time().Before(maxRange.From.Time()) {
			maxRange.From = r.From
		}
		if r.Until.Time().After(maxRange.Until.Time()) {
			maxRange.Until = r.Until
		}
	}

	return maxRange.Ptr()
}

// Min returns the minimum DateRange that is contained within all ranges in the list.
// It finds the latest From date and the earliest Until date across all ranges.
// Returns nil if the list is empty or if there is no overlap between ranges.
func (ranges DateRanges) Min() *DateRange {
	if len(ranges) == 0 {
		return nil
	}

	minRange := ranges[0]
	for _, r := range ranges[1:] {
		if r.From.Time().After(minRange.From.Time()) {
			minRange.From = r.From
		}
		if r.Until.Time().Before(minRange.Until.Time()) {
			minRange.Until = r.Until
		}
	}

	// Check if the resulting range is valid (From <= Until)
	if minRange.From.Time().After(minRange.Until.Time()) {
		return nil
	}

	return minRange.Ptr()
}

// DateRangeFromTime creates a DateRange from two time.Time values.
// It converts the from and until times to Date types and returns a DateRange.
func DateRangeFromTime(from, until stdtime.Time) DateRange {
	return DateRange{
		From:  Date(from),
		Until: Date(until),
	}
}

type DateRange struct {
	From  Date `json:"from,omitempty"`
	Until Date `json:"until,omitempty"`
}

func (r DateRange) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("from", r.From),
		validation.Name("until", r.Until),
		validation.Name("range", validation.HasValidationFunc(func(ctx context.Context) error {
			if r.From.Time().After(r.Until.Time()) {
				return errors.Wrapf(ctx, validation.Error, "from must be less than or equal to until")
			}
			return nil
		})),
	}.Validate(ctx)
}

func (r DateRange) Ptr() *DateRange {
	return &r
}

// DayDateRange creates a DateRange covering the entire day containing the given date.
// The range spans from 00:00:00.000000000 to 23:59:59.999999999 of that day.
func DayDateRange(d Date) DateRange {
	t := d.Time()
	return DateRange{From: Date(BeginningOfDay(t)), Until: Date(EndOfDay(t))}
}

// WeekDateRange creates a DateRange covering the entire week containing the given date.
// The range spans from Monday to Sunday of that week.
// Uses ISO 8601 standard where Monday is the first day of the week.
func WeekDateRange(d Date) DateRange {
	t := d.Time()
	return DateRange{From: Date(BeginningOfWeek(t)), Until: Date(EndOfWeek(t))}
}

// MonthDateRange creates a DateRange covering the entire month containing the given date.
// The range spans from the 1st day to the last day of that month.
func MonthDateRange(d Date) DateRange {
	t := d.Time()
	return DateRange{From: Date(BeginningOfMonth(t)), Until: Date(EndOfMonth(t))}
}

// QuarterDateRange creates a DateRange covering the entire quarter containing the given date.
// Quarters are defined as: Q1=Jan-Mar, Q2=Apr-Jun, Q3=Jul-Sep, Q4=Oct-Dec.
// The range spans from the 1st day of the quarter to the last day of that quarter.
func QuarterDateRange(d Date) DateRange {
	t := d.Time()
	return DateRange{From: Date(BeginningOfQuarter(t)), Until: Date(EndOfQuarter(t))}
}

// YearDateRange creates a DateRange covering the entire year containing the given date.
// The range spans from January 1st to December 31st of that year.
func YearDateRange(d Date) DateRange {
	t := d.Time()
	return DateRange{From: Date(BeginningOfYear(t)), Until: Date(EndOfYear(t))}
}

// TimeRange converts a DateRange to a TimeRange.
func (r DateRange) TimeRange() TimeRange {
	return TimeRange{From: r.From.Time(), Until: r.Until.Time()}
}
