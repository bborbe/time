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

var _ = Describe("FormatTime", func() {
	Describe("with valid time pointer", func() {
		It("formats time in RFC3339 format", func() {
			t := libtimetest.ParseDateTime("2023-12-25T10:15:30Z").Time()

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-12-25T10:15:30Z"))
		})

		It("formats time with timezone offset", func() {
			loc := time.FixedZone("EST", -5*3600)
			t := time.Date(2023, 12, 25, 10, 15, 30, 0, loc)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-12-25T10:15:30-05:00"))
		})

		It("formats time with positive timezone offset", func() {
			loc := time.FixedZone("JST", 9*3600)
			t := time.Date(2023, 12, 25, 10, 15, 30, 0, loc)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-12-25T10:15:30+09:00"))
		})

		It("formats time with nanoseconds", func() {
			t := time.Date(2023, 12, 25, 10, 15, 30, 123456789, time.UTC)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-12-25T10:15:30Z"))
		})

		It("formats zero time", func() {
			t := time.Time{}

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("0001-01-01T00:00:00Z"))
		})

		It("formats time with microsecond precision", func() {
			t := time.Date(2023, 6, 15, 14, 30, 45, 123000000, time.UTC) // 123 milliseconds

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-06-15T14:30:45Z"))
		})

		It("formats time at different hours", func() {
			t := time.Date(2023, 1, 1, 23, 59, 59, 0, time.UTC)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-01-01T23:59:59Z"))
		})

		It("formats time at beginning of day", func() {
			t := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-01-01T00:00:00Z"))
		})
	})

	Describe("with nil pointer", func() {
		It("returns empty string", func() {
			result := libtime.FormatTime(nil)
			Expect(result).To(Equal(""))
		})
	})

	Describe("edge cases", func() {
		It("formats leap year date", func() {
			t := time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC) // Leap year

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2024-02-29T12:00:00Z"))
		})

		It("formats end of year", func() {
			t := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2023-12-31T23:59:59Z"))
		})

		It("formats beginning of year", func() {
			t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2024-01-01T00:00:00Z"))
		})

		It("formats time in different timezones consistently", func() {
			utc := time.UTC
			est := time.FixedZone("EST", -5*3600)
			pst := time.FixedZone("PST", -8*3600)

			// Same moment in different timezones
			baseTime := time.Date(2023, 12, 25, 15, 0, 0, 0, utc)
			estTime := baseTime.In(est)
			pstTime := baseTime.In(pst)

			utcResult := libtime.FormatTime(&baseTime)
			estResult := libtime.FormatTime(&estTime)
			pstResult := libtime.FormatTime(&pstTime)

			Expect(utcResult).To(Equal("2023-12-25T15:00:00Z"))
			Expect(estResult).To(Equal("2023-12-25T10:00:00-05:00"))
			Expect(pstResult).To(Equal("2023-12-25T07:00:00-08:00"))
		})

		It("formats very old dates", func() {
			t := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("1900-01-01T00:00:00Z"))
		})

		It("formats far future dates", func() {
			t := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)

			result := libtime.FormatTime(&t)
			Expect(result).To(Equal("2100-12-31T23:59:59Z"))
		})
	})

	Describe("RFC3339 compliance", func() {
		It("produces parseable RFC3339 strings", func() {
			originalTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)

			formatted := libtime.FormatTime(&originalTime)
			parsedTime, err := time.Parse(time.RFC3339, formatted)

			Expect(err).To(BeNil())
			Expect(parsedTime).To(Equal(originalTime))
		})

		It("produces parseable RFC3339 strings with timezone", func() {
			loc := time.FixedZone("CET", 3600)
			originalTime := time.Date(2023, 6, 15, 14, 30, 45, 0, loc)

			formatted := libtime.FormatTime(&originalTime)
			parsedTime, err := time.Parse(time.RFC3339, formatted)

			Expect(err).To(BeNil())
			// Times should represent the same moment, even if timezone representation differs
			Expect(parsedTime.Equal(originalTime)).To(BeTrue())
		})
	})
})
