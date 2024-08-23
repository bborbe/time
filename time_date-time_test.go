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

var _ = Describe("DateTime", func() {
	var err error
	var snapshotTime libtime.DateTime
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("MarshalJSON", func() {
		var bytes []byte
		BeforeEach(func() {
			snapshotTime = libtime.DateTime(time.Unix(1687161394, 0))
		})
		JustBeforeEach(func() {
			bytes, err = snapshotTime.MarshalJSON()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(string(bytes)).To(Equal(`"2023-06-19T07:56:34Z"`))
		})
	})
	Context("UnmarshalJSON", func() {
		BeforeEach(func() {
			snapshotTime = libtime.DateTime{}
		})
		JustBeforeEach(func() {
			err = snapshotTime.UnmarshalJSON([]byte(`"2023-06-19T07:56:34Z"`))
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(snapshotTime.Time().Format(time.RFC3339Nano)).To(Equal(`2023-06-19T07:56:34Z`))
		})
	})
	Context("ParseDateTime", func() {
		var value interface{}
		var stdTime *libtime.DateTime
		JustBeforeEach(func() {
			stdTime, err = libtime.ParseDateTime(ctx, value)
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
				Expect(stdTime.Format(time.RFC3339)).To(Equal("2023-06-19T07:56:34Z"))
			})
		})
	})
	DescribeTable("ComparePtr",
		func(a *libtime.DateTime, b *libtime.DateTime, expectedResult int) {
			Expect(a.ComparePtr(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.DateTime(time.Unix(1000, 0)).Ptr(), libtime.DateTime(time.Unix(1000, 0)).Ptr(), 0),
		Entry("less", libtime.DateTime(time.Unix(999, 0)).Ptr(), libtime.DateTime(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.DateTime(time.Unix(1000, 0)).Ptr(), libtime.DateTime(time.Unix(999, 0)).Ptr(), 1),
		Entry("equal", nil, nil, 0),
		Entry("less", nil, libtime.DateTime(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.DateTime(time.Unix(1000, 0)).Ptr(), nil, 1),
	)
	DescribeTable("Compare",
		func(a libtime.DateTime, b libtime.DateTime, expectedResult int) {
			Expect(a.Compare(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.DateTime(time.Unix(1000, 0)), libtime.DateTime(time.Unix(1000, 0)), 0),
		Entry("less", libtime.DateTime(time.Unix(999, 0)), libtime.DateTime(time.Unix(1000, 0)), -1),
		Entry("greater", libtime.DateTime(time.Unix(1000, 0)), libtime.DateTime(time.Unix(999, 0)), 1),
	)
})
