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
})
