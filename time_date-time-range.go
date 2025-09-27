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

type DateTimeRanges []DateTimeRange

// Max returns the maximum DateTimeRange that encompasses all ranges in the list.
// It finds the earliest From time and the latest Until time across all ranges.
// Returns nil if the list is empty.
func (ranges DateTimeRanges) Max() *DateTimeRange {
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

// Min returns the minimum DateTimeRange that is contained within all ranges in the list.
// It finds the latest From time and the earliest Until time across all ranges.
// Returns nil if the list is empty or if there is no overlap between ranges.
func (ranges DateTimeRanges) Min() *DateTimeRange {
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

// DateTimeRangeFromTime creates a DateTimeRange from two time.Time values.
// It converts the from and until times to DateTime types and returns a DateTimeRange.
func DateTimeRangeFromTime(from, until stdtime.Time) DateTimeRange {
	return DateTimeRange{
		From:  DateTime(from),
		Until: DateTime(until),
	}
}

type DateTimeRange struct {
	From  DateTime `json:"from,omitempty"`
	Until DateTime `json:"until,omitempty"`
}

func (r DateTimeRange) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("from", r.From),
		validation.Name("until", r.Until),
		validation.Name("range", validation.HasValidationFunc(func(ctx context.Context) error {
			if r.From.After(r.Until) {
				return errors.Wrapf(ctx, validation.Error, "from must be less than or equal to until")
			}
			return nil
		})),
	}.Validate(ctx)
}

func (r DateTimeRange) Ptr() *DateTimeRange {
	return &r
}

// DayDateTimeRange creates a DateTimeRange covering the entire day containing the given datetime.
// The range spans from 00:00:00.000000000 to 23:59:59.999999999 of that day.
func DayDateTimeRange(dt DateTime) DateTimeRange {
	t := dt.Time()
	return DateTimeRange{From: DateTime(BeginningOfDay(t)), Until: DateTime(EndOfDay(t))}
}

// WeekDateTimeRange creates a DateTimeRange covering the entire week containing the given datetime.
// The range spans from Monday 00:00:00.000000000 to Sunday 23:59:59.999999999 of that week.
// Uses ISO 8601 standard where Monday is the first day of the week.
func WeekDateTimeRange(dt DateTime) DateTimeRange {
	t := dt.Time()
	return DateTimeRange{From: DateTime(BeginningOfWeek(t)), Until: DateTime(EndOfWeek(t))}
}

// MonthDateTimeRange creates a DateTimeRange covering the entire month containing the given datetime.
// The range spans from the 1st day 00:00:00.000000000 to the last day 23:59:59.999999999 of that month.
func MonthDateTimeRange(dt DateTime) DateTimeRange {
	t := dt.Time()
	return DateTimeRange{From: DateTime(BeginningOfMonth(t)), Until: DateTime(EndOfMonth(t))}
}

// QuarterDateTimeRange creates a DateTimeRange covering the entire quarter containing the given datetime.
// Quarters are defined as: Q1=Jan-Mar, Q2=Apr-Jun, Q3=Jul-Sep, Q4=Oct-Dec.
// The range spans from the 1st day of the quarter 00:00:00.000000000 to the last day 23:59:59.999999999 of that quarter.
func QuarterDateTimeRange(dt DateTime) DateTimeRange {
	t := dt.Time()
	return DateTimeRange{From: DateTime(BeginningOfQuarter(t)), Until: DateTime(EndOfQuarter(t))}
}

// YearDateTimeRange creates a DateTimeRange covering the entire year containing the given datetime.
// The range spans from January 1st 00:00:00.000000000 to December 31st 23:59:59.999999999 of that year.
func YearDateTimeRange(dt DateTime) DateTimeRange {
	t := dt.Time()
	return DateTimeRange{From: DateTime(BeginningOfYear(t)), Until: DateTime(EndOfYear(t))}
}

// TimeRange converts a DateTimeRange to a TimeRange.
func (r DateTimeRange) TimeRange() TimeRange {
	return TimeRange{From: r.From.Time(), Until: r.Until.Time()}
}
