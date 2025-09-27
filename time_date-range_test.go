// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

var _ = Describe("DateRange", func() {
	Context("DateRangeFromTime", func() {
		var result libtime.DateRange
		var from, until time.Time
		BeforeEach(func() {
			from = time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)
			until = time.Date(2023, time.June, 25, 0, 0, 0, 0, time.UTC)
			result = libtime.DateRangeFromTime(from, until)
		})
		It("creates correct DateRange", func() {
			Expect(result.From.Time()).To(Equal(from))
			Expect(result.Until.Time()).To(Equal(until))
		})
		It("converts time.Time to Date types", func() {
			Expect(result.From.Year()).To(Equal(2023))
			Expect(result.From.Month()).To(Equal(time.June))
			Expect(result.From.Day()).To(Equal(19))
			Expect(result.Until.Year()).To(Equal(2023))
			Expect(result.Until.Month()).To(Equal(time.June))
			Expect(result.Until.Day()).To(Equal(25))
		})
	})

	Context("Range constructors", func() {
		var testDate libtime.Date
		BeforeEach(func() {
			// Tuesday, June 20, 2023
			testDate = libtime.Date(time.Date(2023, time.June, 20, 14, 30, 45, 0, time.UTC))
		})

		Context("DayDateRange", func() {
			It("creates correct day range", func() {
				result := libtime.DayDateRange(testDate)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)))
			})
		})

		Context("WeekDateRange", func() {
			It("creates correct week range", func() {
				result := libtime.WeekDateRange(testDate)
				// Tuesday June 20 -> Week starts Monday June 19, ends Sunday June 25
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)))
			})
		})

		Context("MonthDateRange", func() {
			It("creates correct month range", func() {
				result := libtime.MonthDateRange(testDate)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)))
			})
		})

		Context("QuarterDateRange", func() {
			It("creates correct quarter range", func() {
				result := libtime.QuarterDateRange(testDate)
				// June is Q2 (Apr-Jun)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)))
			})
		})

		Context("YearDateRange", func() {
			It("creates correct year range", func() {
				result := libtime.YearDateRange(testDate)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC)))
			})
		})
	})

	Context("TimeRange conversion", func() {
		It("converts DateRange to TimeRange correctly", func() {
			testDate := libtime.Date(time.Date(2023, time.June, 20, 14, 30, 45, 0, time.UTC))
			dateRange := libtime.DayDateRange(testDate)
			timeRange := dateRange.TimeRange()

			Expect(timeRange.From).To(Equal(dateRange.From.Time()))
			Expect(timeRange.Until).To(Equal(dateRange.Until.Time()))
		})
	})

	Context("Consistency within periods", func() {
		It("produces same day range for all times within the same day", func() {
			dates := []libtime.Date{
				libtime.Date(time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)),            // start
				libtime.Date(time.Date(2023, time.June, 20, 12, 30, 45, 0, time.UTC)),         // midday
				libtime.Date(time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)), // end
			}

			expectedRange := libtime.DayDateRange(dates[0])
			for _, d := range dates {
				result := libtime.DayDateRange(d)
				Expect(result).To(Equal(expectedRange))
			}
		})

		It("produces same month range for all dates within the same month", func() {
			dates := []libtime.Date{
				libtime.Date(time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)),             // first day
				libtime.Date(time.Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC)),           // middle
				libtime.Date(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)), // last day
			}

			expectedRange := libtime.MonthDateRange(dates[0])
			for _, d := range dates {
				result := libtime.MonthDateRange(d)
				Expect(result).To(Equal(expectedRange))
			}
		})
	})
})

var _ = Describe("DateRanges", func() {
	Context("Max", func() {
		It("returns nil for empty list", func() {
			ranges := libtime.DateRanges{}
			result := ranges.Max()
			Expect(result).To(BeNil())
		})

		It("returns the single range for single item list", func() {
			range1 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			)
			ranges := libtime.DateRanges{range1}
			result := ranges.Max()
			Expect(result).NotTo(BeNil())
			Expect(*result).To(Equal(range1))
		})

		It("returns max range encompassing all ranges", func() {
			range1 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			)
			range2 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 5, 0, 0, 0, 0, time.UTC), // Earlier From
				time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC),
			)
			range3 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC), // Later Until
			)
			ranges := libtime.DateRanges{range1, range2, range3}
			result := ranges.Max()

			Expect(result).NotTo(BeNil())
			Expect(result.From.Time()).To(Equal(time.Date(2023, 6, 5, 0, 0, 0, 0, time.UTC)))
			Expect(result.Until.Time()).To(Equal(time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC)))
		})
	})

	Context("Min", func() {
		It("returns nil for empty list", func() {
			ranges := libtime.DateRanges{}
			result := ranges.Min()
			Expect(result).To(BeNil())
		})

		It("returns the single range for single item list", func() {
			range1 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			)
			ranges := libtime.DateRanges{range1}
			result := ranges.Min()
			Expect(result).NotTo(BeNil())
			Expect(*result).To(Equal(range1))
		})

		It("returns min range that overlaps all ranges", func() {
			range1 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC),
			)
			range2 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC), // Later From
				time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC),
			)
			range3 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 18, 0, 0, 0, 0, time.UTC), // Earlier Until
			)
			ranges := libtime.DateRanges{range1, range2, range3}
			result := ranges.Min()

			Expect(result).NotTo(BeNil())
			Expect(result.From.Time()).To(Equal(time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC)))
			Expect(result.Until.Time()).To(Equal(time.Date(2023, 6, 18, 0, 0, 0, 0, time.UTC)))
		})

		It("returns nil when no overlap exists", func() {
			range1 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			)
			range2 := libtime.DateRangeFromTime(
				time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC), // No overlap
				time.Date(2023, 6, 25, 0, 0, 0, 0, time.UTC),
			)
			ranges := libtime.DateRanges{range1, range2}
			result := ranges.Min()

			Expect(result).To(BeNil())
		})
	})
})
