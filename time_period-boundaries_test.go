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

var _ = Describe("Period Boundaries", func() {
	var testTime time.Time
	var testDateTime libtime.DateTime
	var testDate libtime.Date
	var testUnixTime libtime.UnixTime

	BeforeEach(func() {
		// Tuesday, June 20, 2023 14:30:45 UTC
		testTime = time.Date(2023, time.June, 20, 14, 30, 45, 123456789, time.UTC)
		testDateTime = libtime.DateTime(testTime)
		testDate = libtime.Date(testTime)
		testUnixTime = libtime.UnixTime(testTime)
	})

	Context("BeginningOfDay", func() {
		It("works with time.Time", func() {
			result := libtime.BeginningOfDay(testTime)
			expected := time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("works with DateTime via HasTime interface", func() {
			result := libtime.BeginningOfDayFromHasTime(testDateTime)
			expected := time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("works with Date via HasTime interface", func() {
			result := libtime.BeginningOfDayFromHasTime(testDate)
			expected := time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("works with UnixTime via HasTime interface", func() {
			result := libtime.BeginningOfDayFromHasTime(testUnixTime)
			expected := time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("preserves timezone", func() {
			est := time.FixedZone("EST", -5*3600)
			timeInEST := time.Date(2023, time.June, 20, 14, 30, 45, 0, est)
			result := libtime.BeginningOfDay(timeInEST)
			expected := time.Date(2023, time.June, 20, 0, 0, 0, 0, est)
			Expect(result).To(Equal(expected))
		})
	})

	Context("EndOfDay", func() {
		It("works with time.Time", func() {
			result := libtime.EndOfDay(testTime)
			expected := time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("is exactly 1 nanosecond before next day", func() {
			result := libtime.EndOfDay(testTime)
			nextDay := time.Date(2023, time.June, 21, 0, 0, 0, 0, time.UTC)
			Expect(result.Add(time.Nanosecond)).To(Equal(nextDay))
		})
	})

	Context("BeginningOfWeek", func() {
		It("returns Monday for Tuesday input", func() {
			// Tuesday June 20, 2023 -> Monday June 19, 2023
			result := libtime.BeginningOfWeek(testTime)
			expected := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("returns Monday for Sunday input", func() {
			// Sunday June 25, 2023 -> Monday June 19, 2023 (same week)
			sunday := time.Date(2023, time.June, 25, 14, 30, 45, 0, time.UTC)
			result := libtime.BeginningOfWeek(sunday)
			expected := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("returns same Monday for Monday input", func() {
			// Monday June 19, 2023 -> Monday June 19, 2023
			monday := time.Date(2023, time.June, 19, 14, 30, 45, 0, time.UTC)
			result := libtime.BeginningOfWeek(monday)
			expected := time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})
	})

	Context("EndOfWeek", func() {
		It("returns Sunday for Tuesday input", func() {
			// Tuesday June 20, 2023 -> Sunday June 25, 2023
			result := libtime.EndOfWeek(testTime)
			expected := time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("is exactly 1 nanosecond before next week", func() {
			result := libtime.EndOfWeek(testTime)
			nextWeekMonday := time.Date(2023, time.June, 26, 0, 0, 0, 0, time.UTC)
			Expect(result.Add(time.Nanosecond)).To(Equal(nextWeekMonday))
		})
	})

	Context("BeginningOfMonth", func() {
		It("returns first day of month", func() {
			result := libtime.BeginningOfMonth(testTime)
			expected := time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})
	})

	Context("EndOfMonth", func() {
		It("returns last day of month", func() {
			result := libtime.EndOfMonth(testTime)
			expected := time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("handles February in leap year", func() {
			leapYear := time.Date(2024, time.February, 15, 12, 0, 0, 0, time.UTC)
			result := libtime.EndOfMonth(leapYear)
			expected := time.Date(2024, time.February, 29, 23, 59, 59, 999999999, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("handles February in non-leap year", func() {
			nonLeapYear := time.Date(2023, time.February, 15, 12, 0, 0, 0, time.UTC)
			result := libtime.EndOfMonth(nonLeapYear)
			expected := time.Date(2023, time.February, 28, 23, 59, 59, 999999999, time.UTC)
			Expect(result).To(Equal(expected))
		})
	})

	Context("BeginningOfQuarter", func() {
		DescribeTable("returns correct quarter start",
			func(inputMonth time.Month, expectedMonth time.Month) {
				input := time.Date(2023, inputMonth, 15, 12, 0, 0, 0, time.UTC)
				result := libtime.BeginningOfQuarter(input)
				expected := time.Date(2023, expectedMonth, 1, 0, 0, 0, 0, time.UTC)
				Expect(result).To(Equal(expected))
			},
			Entry("Q1 - January", time.January, time.January),
			Entry("Q1 - February", time.February, time.January),
			Entry("Q1 - March", time.March, time.January),
			Entry("Q2 - April", time.April, time.April),
			Entry("Q2 - May", time.May, time.April),
			Entry("Q2 - June", time.June, time.April),
			Entry("Q3 - July", time.July, time.July),
			Entry("Q3 - August", time.August, time.July),
			Entry("Q3 - September", time.September, time.July),
			Entry("Q4 - October", time.October, time.October),
			Entry("Q4 - November", time.November, time.October),
			Entry("Q4 - December", time.December, time.October),
		)
	})

	Context("EndOfQuarter", func() {
		DescribeTable("returns correct quarter end",
			func(inputMonth time.Month, expectedMonth time.Month, expectedDay int) {
				input := time.Date(2023, inputMonth, 15, 12, 0, 0, 0, time.UTC)
				result := libtime.EndOfQuarter(input)
				expected := time.Date(
					2023,
					expectedMonth,
					expectedDay,
					23,
					59,
					59,
					999999999,
					time.UTC,
				)
				Expect(result).To(Equal(expected))
			},
			Entry("Q1 - January", time.January, time.March, 31),
			Entry("Q1 - February", time.February, time.March, 31),
			Entry("Q1 - March", time.March, time.March, 31),
			Entry("Q2 - April", time.April, time.June, 30),
			Entry("Q2 - May", time.May, time.June, 30),
			Entry("Q2 - June", time.June, time.June, 30),
			Entry("Q3 - July", time.July, time.September, 30),
			Entry("Q3 - August", time.August, time.September, 30),
			Entry("Q3 - September", time.September, time.September, 30),
			Entry("Q4 - October", time.October, time.December, 31),
			Entry("Q4 - November", time.November, time.December, 31),
			Entry("Q4 - December", time.December, time.December, 31),
		)
	})

	Context("BeginningOfYear", func() {
		It("returns January 1st", func() {
			result := libtime.BeginningOfYear(testTime)
			expected := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
			Expect(result).To(Equal(expected))
		})
	})

	Context("EndOfYear", func() {
		It("returns December 31st", func() {
			result := libtime.EndOfYear(testTime)
			expected := time.Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC)
			Expect(result).To(Equal(expected))
		})

		It("is exactly 1 nanosecond before next year", func() {
			result := libtime.EndOfYear(testTime)
			nextYear := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
			Expect(result.Add(time.Nanosecond)).To(Equal(nextYear))
		})
	})

	Context("Consistency within periods", func() {
		Context("Day consistency", func() {
			It("produces same boundaries for all times within the same day", func() {
				times := []time.Time{
					time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC),            // start
					time.Date(2023, time.June, 20, 0, 0, 0, 1, time.UTC),            // +1ns
					time.Date(2023, time.June, 20, 12, 30, 45, 123456789, time.UTC), // midday
					time.Date(2023, time.June, 20, 23, 59, 59, 999999998, time.UTC), // almost end
					time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC), // last ns
				}

				expectedStart := libtime.BeginningOfDay(times[0])
				expectedEnd := libtime.EndOfDay(times[0])

				for _, t := range times {
					Expect(libtime.BeginningOfDay(t)).To(Equal(expectedStart))
					Expect(libtime.EndOfDay(t)).To(Equal(expectedEnd))
				}
			})
		})

		Context("Week consistency", func() {
			It("produces same boundaries for all times within the same week", func() {
				// Week June 19-25, 2023 (Monday to Sunday)
				times := []time.Time{
					time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC),            // Monday start
					time.Date(2023, time.June, 20, 14, 30, 0, 0, time.UTC),          // Tuesday
					time.Date(2023, time.June, 22, 9, 15, 0, 0, time.UTC),           // Thursday
					time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC), // Sunday end
				}

				expectedStart := libtime.BeginningOfWeek(times[0])
				expectedEnd := libtime.EndOfWeek(times[0])

				for _, t := range times {
					Expect(libtime.BeginningOfWeek(t)).To(Equal(expectedStart))
					Expect(libtime.EndOfWeek(t)).To(Equal(expectedEnd))
				}
			})
		})

		Context("Month consistency", func() {
			It("produces same boundaries for all times within the same month", func() {
				times := []time.Time{
					time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC),             // start
					time.Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC),           // middle
					time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC), // end
				}

				expectedStart := libtime.BeginningOfMonth(times[0])
				expectedEnd := libtime.EndOfMonth(times[0])

				for _, t := range times {
					Expect(libtime.BeginningOfMonth(t)).To(Equal(expectedStart))
					Expect(libtime.EndOfMonth(t)).To(Equal(expectedEnd))
				}
			})
		})
	})

	Context("Adjacent periods have no gaps", func() {
		It("day boundaries are continuous", func() {
			day1End := libtime.EndOfDay(testTime)
			day2Start := libtime.BeginningOfDay(testTime.AddDate(0, 0, 1))
			Expect(day1End.Add(time.Nanosecond)).To(Equal(day2Start))
		})

		It("month boundaries are continuous", func() {
			monthEnd := libtime.EndOfMonth(testTime)
			nextMonthStart := libtime.BeginningOfMonth(testTime.AddDate(0, 1, 0))
			Expect(monthEnd.Add(time.Nanosecond)).To(Equal(nextMonthStart))
		})

		It("year boundaries are continuous", func() {
			yearEnd := libtime.EndOfYear(testTime)
			nextYearStart := libtime.BeginningOfYear(testTime.AddDate(1, 0, 0))
			Expect(yearEnd.Add(time.Nanosecond)).To(Equal(nextYearStart))
		})
	})
})
