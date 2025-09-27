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

var _ = Describe("DateTimeRange", func() {
	Context("DateTimeRangeFromTime", func() {
		var result libtime.DateTimeRange
		var from, until time.Time
		BeforeEach(func() {
			from = time.Date(2023, time.June, 19, 7, 56, 34, 0, time.UTC)
			until = time.Date(2023, time.June, 25, 15, 30, 45, 0, time.UTC)
			result = libtime.DateTimeRangeFromTime(from, until)
		})
		It("creates correct DateTimeRange", func() {
			Expect(result.From.Time()).To(Equal(from))
			Expect(result.Until.Time()).To(Equal(until))
		})
		It("converts time.Time to DateTime types", func() {
			Expect(result.From.Year()).To(Equal(2023))
			Expect(result.From.Month()).To(Equal(time.June))
			Expect(result.From.Day()).To(Equal(19))
			Expect(result.From.Hour()).To(Equal(7))
			Expect(result.From.Minute()).To(Equal(56))
			Expect(result.From.Second()).To(Equal(34))
			Expect(result.Until.Year()).To(Equal(2023))
			Expect(result.Until.Month()).To(Equal(time.June))
			Expect(result.Until.Day()).To(Equal(25))
			Expect(result.Until.Hour()).To(Equal(15))
			Expect(result.Until.Minute()).To(Equal(30))
			Expect(result.Until.Second()).To(Equal(45))
		})
	})

	Context("Range constructors", func() {
		var testDateTime libtime.DateTime
		BeforeEach(func() {
			// Tuesday, June 20, 2023 14:30:45 UTC
			testDateTime = libtime.DateTime(time.Date(2023, time.June, 20, 14, 30, 45, 123456789, time.UTC))
		})

		Context("DayDateTimeRange", func() {
			It("creates correct day range", func() {
				result := libtime.DayDateTimeRange(testDateTime)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)))
			})

			It("produces same range for all times within the same day", func() {
				dateTimes := []libtime.DateTime{
					libtime.DateTime(time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)),            // start
					libtime.DateTime(time.Date(2023, time.June, 20, 12, 30, 45, 123456789, time.UTC)), // midday
					libtime.DateTime(time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)), // end
				}

				expectedRange := libtime.DayDateTimeRange(dateTimes[0])
				for _, dt := range dateTimes {
					result := libtime.DayDateTimeRange(dt)
					Expect(result).To(Equal(expectedRange))
				}
			})
		})

		Context("WeekDateTimeRange", func() {
			It("creates correct week range", func() {
				result := libtime.WeekDateTimeRange(testDateTime)
				// Tuesday June 20 -> Week starts Monday June 19, ends Sunday June 25
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)))
			})

			It("produces same range for all times within the same week", func() {
				dateTimes := []libtime.DateTime{
					libtime.DateTime(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)),            // Monday start
					libtime.DateTime(time.Date(2023, time.June, 22, 12, 30, 45, 0, time.UTC)),         // Thursday
					libtime.DateTime(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)), // Sunday end
				}

				expectedRange := libtime.WeekDateTimeRange(dateTimes[0])
				for _, dt := range dateTimes {
					result := libtime.WeekDateTimeRange(dt)
					Expect(result).To(Equal(expectedRange))
				}
			})
		})

		Context("MonthDateTimeRange", func() {
			It("creates correct month range", func() {
				result := libtime.MonthDateTimeRange(testDateTime)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)))
			})

			It("handles leap year February", func() {
				leapYear := libtime.DateTime(time.Date(2024, time.February, 15, 12, 0, 0, 0, time.UTC))
				result := libtime.MonthDateTimeRange(leapYear)
				Expect(result.From.Time()).To(Equal(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2024, time.February, 29, 23, 59, 59, 999999999, time.UTC)))
			})
		})

		Context("QuarterDateTimeRange", func() {
			It("creates correct Q2 range for June", func() {
				result := libtime.QuarterDateTimeRange(testDateTime)
				// June is Q2 (Apr-Jun)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)))
			})

			DescribeTable("creates correct quarter range for each month",
				func(inputMonth time.Month, expectedStartMonth time.Month, expectedEndMonth time.Month, expectedEndDay int) {
					input := libtime.DateTime(time.Date(2023, inputMonth, 15, 12, 0, 0, 0, time.UTC))
					result := libtime.QuarterDateTimeRange(input)
					Expect(result.From.Time()).To(Equal(time.Date(2023, expectedStartMonth, 1, 0, 0, 0, 0, time.UTC)))
					Expect(result.Until.Time()).To(Equal(time.Date(2023, expectedEndMonth, expectedEndDay, 23, 59, 59, 999999999, time.UTC)))
				},
				Entry("Q1 - January", time.January, time.January, time.March, 31),
				Entry("Q1 - March", time.March, time.January, time.March, 31),
				Entry("Q2 - April", time.April, time.April, time.June, 30),
				Entry("Q2 - June", time.June, time.April, time.June, 30),
				Entry("Q3 - July", time.July, time.July, time.September, 30),
				Entry("Q3 - September", time.September, time.July, time.September, 30),
				Entry("Q4 - October", time.October, time.October, time.December, 31),
				Entry("Q4 - December", time.December, time.October, time.December, 31),
			)
		})

		Context("YearDateTimeRange", func() {
			It("creates correct year range", func() {
				result := libtime.YearDateTimeRange(testDateTime)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC)))
			})
		})
	})

	Context("TimeRange conversion", func() {
		It("converts DateTimeRange to TimeRange correctly", func() {
			testDateTime := libtime.DateTime(time.Date(2023, time.June, 20, 14, 30, 45, 0, time.UTC))
			dateTimeRange := libtime.DayDateTimeRange(testDateTime)
			timeRange := dateTimeRange.TimeRange()

			Expect(timeRange.From).To(Equal(dateTimeRange.From.Time()))
			Expect(timeRange.Until).To(Equal(dateTimeRange.Until.Time()))
		})

		It("converts all range types to TimeRange consistently", func() {
			testDateTime := libtime.DateTime(time.Date(2023, time.June, 20, 14, 30, 45, 0, time.UTC))

			dayRange := libtime.DayDateTimeRange(testDateTime)
			weekRange := libtime.WeekDateTimeRange(testDateTime)
			monthRange := libtime.MonthDateTimeRange(testDateTime)

			dayTimeRange := dayRange.TimeRange()
			weekTimeRange := weekRange.TimeRange()
			monthTimeRange := monthRange.TimeRange()

			// Verify they all properly convert
			Expect(dayTimeRange.From).To(Equal(dayRange.From.Time()))
			Expect(dayTimeRange.Until).To(Equal(dayRange.Until.Time()))
			Expect(weekTimeRange.From).To(Equal(weekRange.From.Time()))
			Expect(weekTimeRange.Until).To(Equal(weekRange.Until.Time()))
			Expect(monthTimeRange.From).To(Equal(monthRange.From.Time()))
			Expect(monthTimeRange.Until).To(Equal(monthRange.Until.Time()))
		})
	})

	Context("Timezone preservation", func() {
		It("preserves timezone in all range constructors", func() {
			est := time.FixedZone("EST", -5*3600)
			dateTimeInEST := libtime.DateTime(time.Date(2023, time.June, 20, 14, 30, 45, 0, est))

			dayRange := libtime.DayDateTimeRange(dateTimeInEST)
			Expect(dayRange.From.Time().Location()).To(Equal(est))
			Expect(dayRange.Until.Time().Location()).To(Equal(est))

			monthRange := libtime.MonthDateTimeRange(dateTimeInEST)
			Expect(monthRange.From.Time().Location()).To(Equal(est))
			Expect(monthRange.Until.Time().Location()).To(Equal(est))
		})
	})
})

var _ = Describe("DateTimeRanges", func() {
	Context("Max", func() {
		It("returns nil for empty list", func() {
			ranges := libtime.DateTimeRanges{}
			result := ranges.Max()
			Expect(result).To(BeNil())
		})

		It("returns the single range for single item list", func() {
			range1 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 15, 0, 0, 0, time.UTC),
			)
			ranges := libtime.DateTimeRanges{range1}
			result := ranges.Max()
			Expect(result).NotTo(BeNil())
			Expect(*result).To(Equal(range1))
		})

		It("returns max range encompassing all ranges", func() {
			range1 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 15, 0, 0, 0, time.UTC),
			)
			range2 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 5, 8, 0, 0, 0, time.UTC), // Earlier From
				time.Date(2023, 6, 12, 12, 0, 0, 0, time.UTC),
			)
			range3 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 8, 9, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 20, 18, 0, 0, 0, time.UTC), // Later Until
			)
			ranges := libtime.DateTimeRanges{range1, range2, range3}
			result := ranges.Max()

			Expect(result).NotTo(BeNil())
			Expect(result.From.Time()).To(Equal(time.Date(2023, 6, 5, 8, 0, 0, 0, time.UTC)))
			Expect(result.Until.Time()).To(Equal(time.Date(2023, 6, 20, 18, 0, 0, 0, time.UTC)))
		})
	})

	Context("Min", func() {
		It("returns nil for empty list", func() {
			ranges := libtime.DateTimeRanges{}
			result := ranges.Min()
			Expect(result).To(BeNil())
		})

		It("returns the single range for single item list", func() {
			range1 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 15, 0, 0, 0, time.UTC),
			)
			ranges := libtime.DateTimeRanges{range1}
			result := ranges.Min()
			Expect(result).NotTo(BeNil())
			Expect(*result).To(Equal(range1))
		})

		It("returns min range that overlaps all ranges", func() {
			range1 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 20, 15, 0, 0, 0, time.UTC),
			)
			range2 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 12, 8, 0, 0, 0, time.UTC), // Later From
				time.Date(2023, 6, 25, 12, 0, 0, 0, time.UTC),
			)
			range3 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 8, 9, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 18, 18, 0, 0, 0, time.UTC), // Earlier Until
			)
			ranges := libtime.DateTimeRanges{range1, range2, range3}
			result := ranges.Min()

			Expect(result).NotTo(BeNil())
			Expect(result.From.Time()).To(Equal(time.Date(2023, 6, 12, 8, 0, 0, 0, time.UTC)))
			Expect(result.Until.Time()).To(Equal(time.Date(2023, 6, 18, 18, 0, 0, 0, time.UTC)))
		})

		It("returns nil when no overlap exists", func() {
			range1 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 10, 10, 0, 0, 0, time.UTC),
				time.Date(2023, 6, 15, 15, 0, 0, 0, time.UTC),
			)
			range2 := libtime.DateTimeRangeFromTime(
				time.Date(2023, 6, 20, 8, 0, 0, 0, time.UTC), // No overlap
				time.Date(2023, 6, 25, 12, 0, 0, 0, time.UTC),
			)
			ranges := libtime.DateTimeRanges{range1, range2}
			result := ranges.Min()

			Expect(result).To(BeNil())
		})
	})
})
