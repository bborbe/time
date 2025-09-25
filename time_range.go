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

type TimeRange struct {
	From  stdtime.Time `json:"from,omitempty"`
	Until stdtime.Time `json:"until,omitempty"`
}

func (r TimeRange) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("range", validation.HasValidationFunc(func(ctx context.Context) error {
			if r.From.After(r.Until) {
				return errors.Wrapf(ctx, validation.Error, "from must be less than or equal to until")
			}
			return nil
		})),
	}.Validate(ctx)
}

func (r TimeRange) Ptr() *TimeRange {
	return &r
}

// DayTimeRange creates a TimeRange covering the entire day containing the given time.
// The range spans from 00:00:00.000000000 to 23:59:59.999999999 of that day.
func DayTimeRange(t stdtime.Time) TimeRange {
	return TimeRange{From: BeginningOfDay(t), Until: EndOfDay(t)}
}

// WeekTimeRange creates a TimeRange covering the entire week containing the given time.
// The range spans from Monday 00:00:00.000000000 to Sunday 23:59:59.999999999 of that week.
// Uses ISO 8601 standard where Monday is the first day of the week.
func WeekTimeRange(t stdtime.Time) TimeRange {
	return TimeRange{From: BeginningOfWeek(t), Until: EndOfWeek(t)}
}

// MonthTimeRange creates a TimeRange covering the entire month containing the given time.
// The range spans from the 1st day 00:00:00.000000000 to the last day 23:59:59.999999999 of that month.
func MonthTimeRange(t stdtime.Time) TimeRange {
	return TimeRange{From: BeginningOfMonth(t), Until: EndOfMonth(t)}
}

// QuarterTimeRange creates a TimeRange covering the entire quarter containing the given time.
// Quarters are defined as: Q1=Jan-Mar, Q2=Apr-Jun, Q3=Jul-Sep, Q4=Oct-Dec.
// The range spans from the 1st day of the quarter to the last day 23:59:59.999999999 of that quarter.
func QuarterTimeRange(t stdtime.Time) TimeRange {
	return TimeRange{From: BeginningOfQuarter(t), Until: EndOfQuarter(t)}
}

// YearTimeRange creates a TimeRange covering the entire year containing the given time.
// The range spans from January 1st 00:00:00.000000000 to December 31st 23:59:59.999999999 of that year.
func YearTimeRange(t stdtime.Time) TimeRange {
	return TimeRange{From: BeginningOfYear(t), Until: EndOfYear(t)}
}
