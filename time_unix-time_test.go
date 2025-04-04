// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

var _ = Describe("UnixTime", func() {
	var err error
	var snapshotTime libtime.UnixTime
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("MarshalBinary & UnixTimeFromBinary", func() {
		var unixTime libtime.UnixTime
		var binary []byte
		BeforeEach(func() {
			unixTime = libtime.UnixTime(time.Unix(1687161394, 0))
		})
		JustBeforeEach(func() {
			binary, err = unixTime.MarshalBinary()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns binary", func() {
			Expect(binary).NotTo(BeNil())
		})
		It("returns binary", func() {
			unixTimeFromBinary, err := libtime.UnixTimeFromBinary(ctx, binary)
			Expect(err).To(BeNil())
			Expect(unixTimeFromBinary).NotTo(BeNil())
			Expect(unixTimeFromBinary.Unix()).To(Equal(int64(1687161394)))
		})
	})
	Context("MarshalJSON", func() {
		var bytes []byte
		BeforeEach(func() {
			snapshotTime = libtime.UnixTime(time.Unix(1687161394, 0))
		})
		JustBeforeEach(func() {
			bytes, err = snapshotTime.MarshalJSON()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(string(bytes)).To(Equal(`1687161394`))
		})
	})
	Context("UnmarshalJSON", func() {
		BeforeEach(func() {
			snapshotTime = libtime.UnixTime{}
		})
		JustBeforeEach(func() {
			err = snapshotTime.UnmarshalJSON([]byte(`1687161394`))
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(snapshotTime.Time().Format(time.RFC3339Nano)).To(Equal(`2023-06-19T07:56:34Z`))
		})
	})
	DescribeTable("ParseUnixTime",
		func(input any, expectedDateString string, expectedError bool) {
			result, err := libtime.ParseUnixTime(ctx, input)
			if expectedError {
				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Format(time.RFC3339)).To(Equal(expectedDateString))
			}
		},
		Entry("invalid", "banana", "", true),
		Entry("dateTime", "2023-06-19T07:56:34Z", "2023-06-19T07:56:34Z", false),
		Entry("unixTime", 1687161394, "2023-06-19T07:56:34Z", false),
		Entry("unixTimeStr", "1687161394", "2023-06-19T07:56:34Z", false),
	)
	Context("TimePtr", func() {
		var dateTime *libtime.UnixTime
		var timePtr *time.Time
		BeforeEach(func() {
			dateTime = libtime.UnixTime(time.Unix(1000, 0)).Ptr()
		})
		JustBeforeEach(func() {
			timePtr = dateTime.TimePtr()
		})
		Context("datetime not nil", func() {
			It("returns timePtr", func() {
				Expect(timePtr).NotTo(BeNil())
			})
		})
		Context("datetime nil", func() {
			BeforeEach(func() {
				dateTime = nil
			})
			It("returns not timePtr", func() {
				Expect(timePtr).To(BeNil())
			})
		})
	})
})
