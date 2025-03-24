// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	stdtime "time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

var _ = Describe("ParseTime", func() {
	var err error
	var parseTime *stdtime.Time
	var input string
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
		libtime.Now = func() stdtime.Time {
			return stdtime.Unix(1686419205, 0)
		}
	})
	JustBeforeEach(func() {
		parseTime, err = libtime.ParseTime(ctx, input)
	})
	Context("NOW", func() {
		BeforeEach(func() {
			input = "NOW"
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("correct time", func() {
			Expect(parseTime.Unix()).To(Equal(int64(1686419205)))
		})
	})
	Context("NOW-1h", func() {
		BeforeEach(func() {
			input = "NOW-1h"
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("correct time", func() {
			Expect(parseTime.Unix()).To(Equal(int64(1686419205 - 3600)))
		})
	})
	Context("NOW-1d", func() {
		BeforeEach(func() {
			input = "NOW-1d"
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("correct time", func() {
			Expect(parseTime.Unix()).To(Equal(int64(1686419205 - 24*3600)))
		})
	})
	Context("invalid", func() {
		BeforeEach(func() {
			input = "invalid"
		})
		It("returns error", func() {
			Expect(err).NotTo(BeNil())
		})
		It("returns no time", func() {
			Expect(parseTime).To(BeNil())
		})
	})
})
