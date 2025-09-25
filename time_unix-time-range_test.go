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

var _ = Describe("UnixTimeRange", func() {
	Context("UnixTimeRangeFromTime", func() {
		var result libtime.UnixTimeRange
		var from, until time.Time
		BeforeEach(func() {
			from = time.Date(2023, time.June, 19, 7, 56, 34, 0, time.UTC)
			until = time.Date(2023, time.June, 25, 15, 30, 45, 0, time.UTC)
			result = libtime.UnixTimeRangeFromTime(from, until)
		})
		It("creates correct UnixTimeRange", func() {
			Expect(result.From.Time()).To(Equal(from))
			Expect(result.Until.Time()).To(Equal(until))
		})
		It("converts time.Time to UnixTime types", func() {
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
		var testUnixTime libtime.UnixTime
		BeforeEach(func() {
			// Tuesday, June 20, 2023 14:30:45 UTC
			testUnixTime = libtime.UnixTime(time.Date(2023, time.June, 20, 14, 30, 45, 123456789, time.UTC))
		})

		Context("DayUnixTimeRange", func() {
			It("creates correct day range", func() {
				result := libtime.DayUnixTimeRange(testUnixTime)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)))
			})

			It("produces same range for all times within the same day", func() {
				unixTimes := []libtime.UnixTime{
					libtime.UnixTime(time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)),            // start
					libtime.UnixTime(time.Date(2023, time.June, 20, 12, 30, 45, 123456789, time.UTC)), // midday
					libtime.UnixTime(time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)), // end
				}

				expectedRange := libtime.DayUnixTimeRange(unixTimes[0])
				for _, ut := range unixTimes {
					result := libtime.DayUnixTimeRange(ut)
					Expect(result).To(Equal(expectedRange))
				}
			})
		})

		Context("WeekUnixTimeRange", func() {
			It("creates correct week range", func() {
				result := libtime.WeekUnixTimeRange(testUnixTime)
				// Tuesday June 20 -> Week starts Monday June 19, ends Sunday June 25
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)))
			})

			It("produces same range for all times within the same week", func() {
				unixTimes := []libtime.UnixTime{
					libtime.UnixTime(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)),            // Monday start
					libtime.UnixTime(time.Date(2023, time.June, 22, 12, 30, 45, 0, time.UTC)),         // Thursday
					libtime.UnixTime(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)), // Sunday end
				}

				expectedRange := libtime.WeekUnixTimeRange(unixTimes[0])
				for _, ut := range unixTimes {
					result := libtime.WeekUnixTimeRange(ut)
					Expect(result).To(Equal(expectedRange))
				}
			})
		})

		Context("MonthUnixTimeRange", func() {
			It("creates correct month range", func() {
				result := libtime.MonthUnixTimeRange(testUnixTime)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)))
			})

			It("handles leap year February", func() {
				leapYear := libtime.UnixTime(time.Date(2024, time.February, 15, 12, 0, 0, 0, time.UTC))
				result := libtime.MonthUnixTimeRange(leapYear)
				Expect(result.From.Time()).To(Equal(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2024, time.February, 29, 23, 59, 59, 999999999, time.UTC)))
			})
		})

		Context("QuarterUnixTimeRange", func() {
			It("creates correct Q2 range for June", func() {
				result := libtime.QuarterUnixTimeRange(testUnixTime)
				// June is Q2 (Apr-Jun)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)))
			})

			DescribeTable("creates correct quarter range for each month",
				func(inputMonth time.Month, expectedStartMonth time.Month, expectedEndMonth time.Month, expectedEndDay int) {
					input := libtime.UnixTime(time.Date(2023, inputMonth, 15, 12, 0, 0, 0, time.UTC))
					result := libtime.QuarterUnixTimeRange(input)
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

		Context("YearUnixTimeRange", func() {
			It("creates correct year range", func() {
				result := libtime.YearUnixTimeRange(testUnixTime)
				Expect(result.From.Time()).To(Equal(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)))
				Expect(result.Until.Time()).To(Equal(time.Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC)))
			})
		})
	})

	Context("TimeRange conversion", func() {
		It("converts UnixTimeRange to TimeRange correctly", func() {
			testUnixTime := libtime.UnixTime(time.Date(2023, time.June, 20, 14, 30, 45, 0, time.UTC))
			unixTimeRange := libtime.DayUnixTimeRange(testUnixTime)
			timeRange := unixTimeRange.TimeRange()

			Expect(timeRange.From).To(Equal(unixTimeRange.From.Time()))
			Expect(timeRange.Until).To(Equal(unixTimeRange.Until.Time()))
		})

		It("converts all range types to TimeRange consistently", func() {
			testUnixTime := libtime.UnixTime(time.Date(2023, time.June, 20, 14, 30, 45, 0, time.UTC))

			dayRange := libtime.DayUnixTimeRange(testUnixTime)
			weekRange := libtime.WeekUnixTimeRange(testUnixTime)
			monthRange := libtime.MonthUnixTimeRange(testUnixTime)

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
			unixTimeInEST := libtime.UnixTime(time.Date(2023, time.June, 20, 14, 30, 45, 0, est))

			dayRange := libtime.DayUnixTimeRange(unixTimeInEST)
			Expect(dayRange.From.Time().Location()).To(Equal(est))
			Expect(dayRange.Until.Time().Location()).To(Equal(est))

			monthRange := libtime.MonthUnixTimeRange(unixTimeInEST)
			Expect(monthRange.From.Time().Location()).To(Equal(est))
			Expect(monthRange.Until.Time().Location()).To(Equal(est))
		})
	})
})
