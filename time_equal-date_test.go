// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
	libtimetest "github.com/bborbe/time/test"
)

var _ = Describe("HasEqualDate", func() {
	Describe("with same dates", func() {
		It("returns true for identical times", func() {
			t1 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()
			t2 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})

		It("returns true for same date but different times", func() {
			t1 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()
			t2 := libtimetest.ParseDateTime("2023-12-25T23:59:59Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})

		It("returns true for same date with different seconds", func() {
			t1 := libtimetest.ParseDateTime("2023-06-15T08:30:15Z").Time()
			t2 := libtimetest.ParseDateTime("2023-06-15T08:30:45Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})

		It("returns true for same date with different nanoseconds", func() {
			t1 := time.Date(2023, 3, 10, 12, 0, 0, 123456789, time.UTC)
			t2 := time.Date(2023, 3, 10, 12, 0, 0, 987654321, time.UTC)

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})
	})

	Describe("with different dates", func() {
		It("returns false for different years", func() {
			t1 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()
			t2 := libtimetest.ParseDateTime("2024-12-25T10:15:30Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeFalse())
		})

		It("returns false for different months", func() {
			t1 := libtimetest.ParseDateTime("2023-11-25T10:15:30Z").Time()
			t2 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeFalse())
		})

		It("returns false for different days", func() {
			t1 := libtimetest.ParseDateTime("2023-12-24T10:15:30Z").Time()
			t2 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeFalse())
		})

		It("returns false for completely different dates", func() {
			t1 := libtimetest.ParseDateTime("2020-01-01T00:00:00Z").Time()
			t2 := libtimetest.ParseDateTime("2025-12-31T23:59:59Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeFalse())
		})
	})

	Describe("with different timezones", func() {
		It("returns true when same date in different timezones", func() {
			// These represent the same moment but in different timezones
			// Both are 2023-12-25 in their respective zones
			t1 := time.Date(2023, 12, 25, 10, 0, 0, 0, time.UTC)
			t2 := time.Date(2023, 12, 25, 15, 0, 0, 0, time.FixedZone("EST", -5*3600))

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})

		It("returns false when different dates due to timezone", func() {
			// These represent the same moment but different calendar dates
			// t1 is 2023-12-25 UTC, t2 is 2023-12-24 in its zone
			utc := time.UTC
			pst := time.FixedZone("PST", -8*3600)

			t1 := time.Date(2023, 12, 25, 2, 0, 0, 0, utc)  // 2023-12-25 02:00 UTC
			t2 := time.Date(2023, 12, 24, 18, 0, 0, 0, pst) // 2023-12-24 18:00 PST (same moment as above)

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeFalse())
		})

		It("compares dates in each time's respective timezone", func() {
			// The function compares dates as they appear in each timezone, not UTC
			est := time.FixedZone("EST", -5*3600)
			pst := time.FixedZone("PST", -8*3600)

			// Both are December 25th in their respective timezones
			t1 := time.Date(2023, 12, 25, 10, 0, 0, 0, est) // 2023-12-25 10:00 EST
			t2 := time.Date(2023, 12, 25, 7, 0, 0, 0, pst)  // 2023-12-25 07:00 PST

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})
	})

	Describe("edge cases", func() {
		It("returns true for zero times", func() {
			t1 := time.Time{}
			t2 := time.Time{}

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})

		It("returns false when one time is zero", func() {
			t1 := time.Time{}
			t2 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeFalse())
		})

		It("handles leap year dates correctly", func() {
			// February 29th in leap year vs non-leap year
			leapYear := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)
			nonLeapYear := time.Date(2023, 2, 28, 12, 0, 0, 0, time.UTC) // 2023 doesn't have Feb 29

			result := libtime.HasEqualDate(leapYear, nonLeapYear)
			Expect(result).To(BeFalse())
		})

		It("handles end of year/beginning of year boundary", func() {
			endOfYear := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
			beginningOfYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

			result := libtime.HasEqualDate(endOfYear, beginningOfYear)
			Expect(result).To(BeFalse())
		})

		It("handles same date on year boundary", func() {
			t1 := time.Date(2023, 12, 31, 10, 0, 0, 0, time.UTC)
			t2 := time.Date(2023, 12, 31, 22, 0, 0, 0, time.UTC)

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue())
		})

		It("handles daylight saving time transitions", func() {
			// Create times around DST transition (this is conceptual since actual DST handling depends on location)
			loc, err := time.LoadLocation("America/New_York")
			Expect(err).To(BeNil())

			// Both are March 12th, 2023 in New York, but one might be during DST transition
			t1 := time.Date(2023, 3, 12, 1, 0, 0, 0, loc) // Before DST
			t2 := time.Date(2023, 3, 12, 3, 0, 0, 0, loc) // After DST (2 AM gets skipped)

			result := libtime.HasEqualDate(t1, t2)
			Expect(result).To(BeTrue()) // Same calendar date despite DST
		})
	})
})
