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

var _ = Describe("TimeRange Constructors", func() {
	var testTime time.Time
	BeforeEach(func() {
		// Tuesday, June 20, 2023 14:30:45 UTC
		testTime = time.Date(2023, time.June, 20, 14, 30, 45, 123456789, time.UTC)
	})

	Context("DayTimeRange", func() {
		It("creates correct day range", func() {
			result := libtime.DayTimeRange(testTime)
			Expect(result.From).To(Equal(time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC)))
			Expect(
				result.Until,
			).To(Equal(time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC)))
		})

		It("produces same range for all times within the same day", func() {
			times := []time.Time{
				time.Date(2023, time.June, 20, 0, 0, 0, 0, time.UTC),            // start
				time.Date(2023, time.June, 20, 12, 30, 45, 123456789, time.UTC), // midday
				time.Date(2023, time.June, 20, 23, 59, 59, 999999999, time.UTC), // end
			}

			expectedRange := libtime.DayTimeRange(times[0])
			for _, t := range times {
				result := libtime.DayTimeRange(t)
				Expect(result).To(Equal(expectedRange))
			}
		})
	})

	Context("WeekTimeRange", func() {
		It("creates correct week range", func() {
			result := libtime.WeekTimeRange(testTime)
			// Tuesday June 20 -> Week starts Monday June 19, ends Sunday June 25
			Expect(result.From).To(Equal(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)))
			Expect(
				result.Until,
			).To(Equal(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)))
		})

		It("handles Sunday correctly", func() {
			sunday := time.Date(2023, time.June, 25, 14, 30, 45, 0, time.UTC)
			result := libtime.WeekTimeRange(sunday)
			// Sunday June 25 -> Same week (June 19-25)
			Expect(result.From).To(Equal(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)))
			Expect(
				result.Until,
			).To(Equal(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)))
		})

		It("handles Monday correctly", func() {
			monday := time.Date(2023, time.June, 19, 14, 30, 45, 0, time.UTC)
			result := libtime.WeekTimeRange(monday)
			// Monday June 19 -> Same week (June 19-25)
			Expect(result.From).To(Equal(time.Date(2023, time.June, 19, 0, 0, 0, 0, time.UTC)))
			Expect(
				result.Until,
			).To(Equal(time.Date(2023, time.June, 25, 23, 59, 59, 999999999, time.UTC)))
		})
	})

	Context("MonthTimeRange", func() {
		It("creates correct month range", func() {
			result := libtime.MonthTimeRange(testTime)
			Expect(result.From).To(Equal(time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)))
			Expect(
				result.Until,
			).To(Equal(time.Date(2023, time.June, 30, 23, 59, 59, 999999999, time.UTC)))
		})

		It("handles leap year February", func() {
			leapYear := time.Date(2024, time.February, 15, 12, 0, 0, 0, time.UTC)
			result := libtime.MonthTimeRange(leapYear)
			Expect(result.From).To(Equal(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)))
			Expect(
				result.Until,
			).To(Equal(time.Date(2024, time.February, 29, 23, 59, 59, 999999999, time.UTC)))
		})
	})

	Context("QuarterTimeRange", func() {
		DescribeTable(
			"creates correct quarter range",
			func(inputMonth time.Month, expectedStartMonth time.Month, expectedEndMonth time.Month, expectedEndDay int) {
				input := time.Date(2023, inputMonth, 15, 12, 0, 0, 0, time.UTC)
				result := libtime.QuarterTimeRange(input)
				Expect(
					result.From,
				).To(Equal(time.Date(2023, expectedStartMonth, 1, 0, 0, 0, 0, time.UTC)))
				Expect(
					result.Until,
				).To(Equal(time.Date(2023, expectedEndMonth, expectedEndDay, 23, 59, 59, 999999999, time.UTC)))
			},
			Entry("Q1 - January", time.January, time.January, time.March, 31),
			Entry("Q1 - February", time.February, time.January, time.March, 31),
			Entry("Q1 - March", time.March, time.January, time.March, 31),
			Entry("Q2 - April", time.April, time.April, time.June, 30),
			Entry("Q2 - May", time.May, time.April, time.June, 30),
			Entry("Q2 - June", time.June, time.April, time.June, 30),
			Entry("Q3 - July", time.July, time.July, time.September, 30),
			Entry("Q3 - August", time.August, time.July, time.September, 30),
			Entry("Q3 - September", time.September, time.July, time.September, 30),
			Entry("Q4 - October", time.October, time.October, time.December, 31),
			Entry("Q4 - November", time.November, time.October, time.December, 31),
			Entry("Q4 - December", time.December, time.October, time.December, 31),
		)
	})

	Context("YearTimeRange", func() {
		It("creates correct year range", func() {
			result := libtime.YearTimeRange(testTime)
			Expect(result.From).To(Equal(time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)))
			Expect(
				result.Until,
			).To(Equal(time.Date(2023, time.December, 31, 23, 59, 59, 999999999, time.UTC)))
		})
	})

	Context("Timezone preservation", func() {
		It("preserves timezone in all range constructors", func() {
			est := time.FixedZone("EST", -5*3600)
			timeInEST := time.Date(2023, time.June, 20, 14, 30, 45, 0, est)

			dayRange := libtime.DayTimeRange(timeInEST)
			Expect(dayRange.From.Location()).To(Equal(est))
			Expect(dayRange.Until.Location()).To(Equal(est))

			weekRange := libtime.WeekTimeRange(timeInEST)
			Expect(weekRange.From.Location()).To(Equal(est))
			Expect(weekRange.Until.Location()).To(Equal(est))

			monthRange := libtime.MonthTimeRange(timeInEST)
			Expect(monthRange.From.Location()).To(Equal(est))
			Expect(monthRange.Until.Location()).To(Equal(est))
		})
	})
})
