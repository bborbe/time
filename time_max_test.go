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

var _ = Describe("Max", func() {
	Describe("with different times", func() {
		It("returns the later time when first is earlier", func() {
			earlier := libtimetest.ParseDateTime("2023-12-25T10:00:00Z").Time()
			later := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})

		It("returns the later time when second is earlier", func() {
			earlier := libtimetest.ParseDateTime("2023-12-25T10:00:00Z").Time()
			later := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.Max(later, earlier)
			Expect(result).To(Equal(later))
		})

		It("returns the later time with significant time difference", func() {
			earlier := libtimetest.ParseDateTime("2020-01-01T00:00:00Z").Time()
			later := libtimetest.ParseDateTime("2025-12-31T23:59:59Z").Time()

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})

		It("returns the later time with different dates", func() {
			earlier := libtimetest.ParseDateTime("2023-12-24T23:59:59Z").Time()
			later := libtimetest.ParseDateTime("2023-12-25T00:00:00Z").Time()

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})

		It("returns the later time with different years", func() {
			earlier := libtimetest.ParseDateTime("2022-12-31T23:59:59Z").Time()
			later := libtimetest.ParseDateTime("2023-01-01T00:00:00Z").Time()

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})
	})

	Describe("with equal times", func() {
		It("returns the first time when both are identical", func() {
			t1 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()
			t2 := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.Max(t1, t2)
			Expect(result).To(Equal(t1))
			Expect(result).To(Equal(t2))
		})

		It("returns the first time when both represent the same moment", func() {
			t1 := time.Date(2023, 12, 25, 10, 15, 30, 123456789, time.UTC)
			t2 := time.Date(2023, 12, 25, 10, 15, 30, 123456789, time.UTC)

			result := libtime.Max(t1, t2)
			Expect(result).To(Equal(t1))
		})
	})

	Describe("with nanosecond precision", func() {
		It("returns the later time with nanosecond difference", func() {
			earlier := time.Date(2023, 12, 25, 10, 15, 30, 123456789, time.UTC)
			later := time.Date(2023, 12, 25, 10, 15, 30, 123456790, time.UTC)

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})

		It("returns the later time with microsecond difference", func() {
			earlier := time.Date(2023, 12, 25, 10, 15, 30, 123456000, time.UTC)
			later := time.Date(2023, 12, 25, 10, 15, 30, 123457000, time.UTC)

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})

		It("returns the later time with millisecond difference", func() {
			earlier := time.Date(2023, 12, 25, 10, 15, 30, 123000000, time.UTC)
			later := time.Date(2023, 12, 25, 10, 15, 30, 124000000, time.UTC)

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})
	})

	Describe("with different timezones", func() {
		It("compares times correctly across timezones", func() {
			utc := time.UTC
			est := time.FixedZone("EST", -5*3600)

			// Same moment in different timezones
			t1 := time.Date(2023, 12, 25, 15, 0, 0, 0, utc) // 15:00 UTC
			t2 := time.Date(2023, 12, 25, 10, 0, 0, 0, est) // 10:00 EST (same moment)

			result := libtime.Max(t1, t2)
			// Should return one of them since they're the same moment
			Expect(result.Equal(t1)).To(BeTrue())
			Expect(result.Equal(t2)).To(BeTrue())
		})

		It("returns the later time when times are in different zones", func() {
			utc := time.UTC
			jst := time.FixedZone("JST", 9*3600)

			earlier := time.Date(2023, 12, 25, 10, 0, 0, 0, utc) // 10:00 UTC
			later := time.Date(2023, 12, 25, 20, 0, 0, 0, jst)   // 20:00 JST (11:00 UTC)

			result := libtime.Max(earlier, later)
			Expect(result).To(Equal(later))
		})

		It("handles timezone comparisons correctly", func() {
			pst := time.FixedZone("PST", -8*3600)
			est := time.FixedZone("EST", -5*3600)

			// Different moments in different timezones
			pstTime := time.Date(2023, 12, 25, 10, 0, 0, 0, pst) // 10:00 PST (18:00 UTC)
			estTime := time.Date(2023, 12, 25, 14, 0, 0, 0, est) // 14:00 EST (19:00 UTC)

			result := libtime.Max(pstTime, estTime)
			Expect(result).To(Equal(estTime)) // EST time is later
		})
	})

	Describe("edge cases", func() {
		It("handles zero times", func() {
			zeroTime := time.Time{}
			nonZeroTime := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.Max(zeroTime, nonZeroTime)
			Expect(result).To(Equal(nonZeroTime))
		})

		It("handles both zero times", func() {
			zeroTime1 := time.Time{}
			zeroTime2 := time.Time{}

			result := libtime.Max(zeroTime1, zeroTime2)
			Expect(result).To(Equal(zeroTime1))
			Expect(result).To(Equal(zeroTime2))
		})

		It("handles zero time as second parameter", func() {
			nonZeroTime := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()
			zeroTime := time.Time{}

			result := libtime.Max(nonZeroTime, zeroTime)
			Expect(result).To(Equal(nonZeroTime))
		})

		It("handles very old dates", func() {
			veryOld := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
			old := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)

			result := libtime.Max(veryOld, old)
			Expect(result).To(Equal(old))
		})

		It("handles far future dates", func() {
			future := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
			farFuture := time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC)

			result := libtime.Max(future, farFuture)
			Expect(result).To(Equal(farFuture))
		})

		It("handles leap year dates", func() {
			leapDay := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC)
			dayAfter := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)

			result := libtime.Max(leapDay, dayAfter)
			Expect(result).To(Equal(dayAfter))
		})

		It("handles daylight saving time transitions", func() {
			loc, err := time.LoadLocation("America/New_York")
			Expect(err).To(BeNil())

			// Before DST transition
			beforeDST := time.Date(2023, 3, 12, 1, 30, 0, 0, loc)
			// After DST transition (2 AM becomes 3 AM)
			afterDST := time.Date(2023, 3, 12, 3, 30, 0, 0, loc)

			result := libtime.Max(beforeDST, afterDST)
			Expect(result).To(Equal(afterDST))
		})
	})

	Describe("symmetry property", func() {
		It("is commutative", func() {
			t1 := libtimetest.ParseDateTime("2023-12-25T10:00:00Z").Time()
			t2 := libtimetest.ParseDateTime("2023-12-25T11:00:00Z").Time()

			result1 := libtime.Max(t1, t2)
			result2 := libtime.Max(t2, t1)

			Expect(result1).To(Equal(result2))
			Expect(result1).To(Equal(t2)) // t2 is later
		})

		It("is idempotent", func() {
			t := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.Max(t, t)
			Expect(result).To(Equal(t))
		})
	})
})
