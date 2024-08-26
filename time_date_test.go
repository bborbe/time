// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"time"

	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Date", func() {
	var err error
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("MarshalJSON", func() {
		var snapshotTime libtime.Date
		var bytes []byte
		JustBeforeEach(func() {
			bytes, err = snapshotTime.MarshalJSON()
		})
		Context("defined", func() {
			BeforeEach(func() {
				snapshotTime = libtime.Date(time.Unix(1687161394, 0))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(string(bytes)).To(Equal(`"2023-06-19"`))
			})
		})
		Context("undefined", func() {
			BeforeEach(func() {
				snapshotTime = libtime.Date{}
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(string(bytes)).To(Equal(`null`))
			})
		})
	})
	Context("UnmarshalJSON", func() {
		var snapshotTime libtime.Date
		var value string
		BeforeEach(func() {
			snapshotTime = libtime.Date{}
		})
		JustBeforeEach(func() {
			err = snapshotTime.UnmarshalJSON([]byte(value))
		})
		Context("with value", func() {
			BeforeEach(func() {
				value = `"2023-06-19"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(snapshotTime.Time().Format(time.DateOnly)).To(Equal(`2023-06-19`))
			})
		})
		Context("with empty value", func() {
			BeforeEach(func() {
				value = `""`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(snapshotTime.Time().IsZero()).To(BeTrue())
			})
		})
		Context("with null value", func() {
			BeforeEach(func() {
				value = `null`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(snapshotTime.Time().IsZero()).To(BeTrue())
			})
		})
	})
	Context("ParseDate", func() {
		var value interface{}
		var stdTime *libtime.Date
		JustBeforeEach(func() {
			stdTime, err = libtime.ParseDate(ctx, value)
		})
		Context("Success", func() {
			BeforeEach(func() {
				value = "2023-06-19T07:56:34Z"
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct time", func() {
				Expect(stdTime).NotTo(BeNil())
				Expect(stdTime.Format(time.DateOnly)).To(Equal("2023-06-19"))
			})
		})
	})
	DescribeTable("ComparePtr",
		func(a *libtime.Date, b *libtime.Date, expectedResult int) {
			Expect(a.ComparePtr(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.Date(time.Unix(1000, 0)).Ptr(), libtime.Date(time.Unix(1000, 0)).Ptr(), 0),
		Entry("less", libtime.Date(time.Unix(999, 0)).Ptr(), libtime.Date(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.Date(time.Unix(1000, 0)).Ptr(), libtime.Date(time.Unix(999, 0)).Ptr(), 1),
		Entry("equal", nil, nil, 0),
		Entry("less", nil, libtime.Date(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.Date(time.Unix(1000, 0)).Ptr(), nil, 1),
	)
	DescribeTable("Compare",
		func(a libtime.Date, b libtime.Date, expectedResult int) {
			Expect(a.Compare(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.Date(time.Unix(1000, 0)), libtime.Date(time.Unix(1000, 0)), 0),
		Entry("less", libtime.Date(time.Unix(999, 0)), libtime.Date(time.Unix(1000, 0)), -1),
		Entry("greater", libtime.Date(time.Unix(1000, 0)), libtime.Date(time.Unix(999, 0)), 1),
	)
})
