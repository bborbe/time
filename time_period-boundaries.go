// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	stdtime "time"
)

// BeginningOfDay returns the start of the day (00:00:00.000000000) for the given time.
// Preserves the original timezone.
func BeginningOfDay(t stdtime.Time) stdtime.Time {
	return stdtime.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// BeginningOfDayFromHasTime returns the start of the day for any type implementing HasTime interface.
func BeginningOfDayFromHasTime(hasTime HasTime) stdtime.Time {
	return BeginningOfDay(hasTime.Time())
}

// EndOfDay returns the end of the day (23:59:59.999999999) for the given time.
// This is calculated as the start of the next day minus 1 nanosecond.
func EndOfDay(t stdtime.Time) stdtime.Time {
	return BeginningOfDay(t).AddDate(0, 0, 1).Add(-stdtime.Nanosecond)
}

// EndOfDayFromHasTime returns the end of the day for any type implementing HasTime interface.
func EndOfDayFromHasTime(hasTime HasTime) stdtime.Time {
	return EndOfDay(hasTime.Time())
}

// BeginningOfWeek returns the start of the week (Monday 00:00:00.000000000) for the given time.
// Uses ISO 8601 standard where Monday is the first day of the week.
// Preserves the original timezone.
func BeginningOfWeek(t stdtime.Time) stdtime.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Go treats Sunday as 0, but ISO 8601 treats it as 7
		weekday = 7
	}
	return stdtime.Date(t.Year(), t.Month(), t.Day()-weekday+1, 0, 0, 0, 0, t.Location())
}

// BeginningOfWeekFromHasTime returns the start of the week for any type implementing HasTime interface.
func BeginningOfWeekFromHasTime(hasTime HasTime) stdtime.Time {
	return BeginningOfWeek(hasTime.Time())
}

// EndOfWeek returns the end of the week (Sunday 23:59:59.999999999) for the given time.
// This is calculated as the start of the next week minus 1 nanosecond.
func EndOfWeek(t stdtime.Time) stdtime.Time {
	return BeginningOfWeek(t).AddDate(0, 0, 7).Add(-stdtime.Nanosecond)
}

// EndOfWeekFromHasTime returns the end of the week for any type implementing HasTime interface.
func EndOfWeekFromHasTime(hasTime HasTime) stdtime.Time {
	return EndOfWeek(hasTime.Time())
}

// BeginningOfMonth returns the start of the month (1st day 00:00:00.000000000) for the given time.
// Preserves the original timezone.
func BeginningOfMonth(t stdtime.Time) stdtime.Time {
	return stdtime.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// BeginningOfMonthFromHasTime returns the start of the month for any type implementing HasTime interface.
func BeginningOfMonthFromHasTime(hasTime HasTime) stdtime.Time {
	return BeginningOfMonth(hasTime.Time())
}

// EndOfMonth returns the end of the month (last day 23:59:59.999999999) for the given time.
// This is calculated as the start of the next month minus 1 nanosecond.
func EndOfMonth(t stdtime.Time) stdtime.Time {
	return BeginningOfMonth(t).AddDate(0, 1, 0).Add(-stdtime.Nanosecond)
}

// EndOfMonthFromHasTime returns the end of the month for any type implementing HasTime interface.
func EndOfMonthFromHasTime(hasTime HasTime) stdtime.Time {
	return EndOfMonth(hasTime.Time())
}

// BeginningOfQuarter returns the start of the quarter (1st day of quarter 00:00:00.000000000) for the given time.
// Quarters are: Q1=Jan-Mar, Q2=Apr-Jun, Q3=Jul-Sep, Q4=Oct-Dec.
// Preserves the original timezone.
func BeginningOfQuarter(t stdtime.Time) stdtime.Time {
	quarterMonth := ((int(t.Month())-1)/3)*3 + 1
	return stdtime.Date(t.Year(), stdtime.Month(quarterMonth), 1, 0, 0, 0, 0, t.Location())
}

// BeginningOfQuarterFromHasTime returns the start of the quarter for any type implementing HasTime interface.
func BeginningOfQuarterFromHasTime(hasTime HasTime) stdtime.Time {
	return BeginningOfQuarter(hasTime.Time())
}

// EndOfQuarter returns the end of the quarter (last day of quarter 23:59:59.999999999) for the given time.
// This is calculated as the start of the next quarter minus 1 nanosecond.
func EndOfQuarter(t stdtime.Time) stdtime.Time {
	return BeginningOfQuarter(t).AddDate(0, 3, 0).Add(-stdtime.Nanosecond)
}

// EndOfQuarterFromHasTime returns the end of the quarter for any type implementing HasTime interface.
func EndOfQuarterFromHasTime(hasTime HasTime) stdtime.Time {
	return EndOfQuarter(hasTime.Time())
}

// BeginningOfYear returns the start of the year (January 1st 00:00:00.000000000) for the given time.
// Preserves the original timezone.
func BeginningOfYear(t stdtime.Time) stdtime.Time {
	return stdtime.Date(t.Year(), stdtime.January, 1, 0, 0, 0, 0, t.Location())
}

// BeginningOfYearFromHasTime returns the start of the year for any type implementing HasTime interface.
func BeginningOfYearFromHasTime(hasTime HasTime) stdtime.Time {
	return BeginningOfYear(hasTime.Time())
}

// EndOfYear returns the end of the year (December 31st 23:59:59.999999999) for the given time.
// This is calculated as the start of the next year minus 1 nanosecond.
func EndOfYear(t stdtime.Time) stdtime.Time {
	return BeginningOfYear(t).AddDate(1, 0, 0).Add(-stdtime.Nanosecond)
}

// EndOfYearFromHasTime returns the end of the year for any type implementing HasTime interface.
func EndOfYearFromHasTime(hasTime HasTime) stdtime.Time {
	return EndOfYear(hasTime.Time())
}
