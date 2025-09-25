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
})
