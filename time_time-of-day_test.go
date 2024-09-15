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

var _ = Describe("TimeOfDay", func() {
	var err error
	var timeOfDay libtime.TimeOfDay
	var now time.Time
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
		now = ParseTime("2023-05-02T12:45:59.123456Z")
	})
	JustBeforeEach(func() {
		libtime.Now = func() time.Time {
			return now
		}
	})
	Context("ParseTimeOfDay", func() {
		var input string
		var timeOfDay *libtime.TimeOfDay
		var winterTime *libtime.DateTime
		var summerTime *libtime.DateTime
		JustBeforeEach(func() {
			timeOfDay, err = libtime.ParseTimeOfDay(ctx, input)
			Expect(err).To(BeNil())
			{
				winterTime, err = timeOfDay.DateTime(2024, time.January, 1)
				Expect(err).To(BeNil())
			}
			{
				summerTime, err = timeOfDay.DateTime(2024, time.July, 1)
				Expect(err).To(BeNil())
			}
		})
		Context("NOW", func() {
			BeforeEach(func() {
				input = "NOW"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T12:45:59.123456Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T12:45:59.123456Z"))
			})
		})
		Context("time with Z", func() {
			BeforeEach(func() {
				input = "13:37:59.123456Z"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T13:37:59.123456Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T13:37:59.123456Z"))
			})
		})
		Context("time with UTC", func() {
			BeforeEach(func() {
				input = "14:37:59 UTC"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T14:37:59Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T14:37:59Z"))
			})
		})
		Context("time with Europe/Berlin", func() {
			BeforeEach(func() {
				input = "15:37:59 Europe/Berlin"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T15:37:59+01:00"))
				Expect(winterTime.Time().UTC().Format(time.RFC3339Nano)).To(Equal("2024-01-01T14:37:59Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T15:37:59+02:00"))
				Expect(summerTime.Time().UTC().Format(time.RFC3339Nano)).To(Equal("2024-07-01T13:37:59Z"))
			})
		})
	})
	Context("String", func() {
		var result string
		JustBeforeEach(func() {
			result = timeOfDay.String()
		})
		Context("with nano", func() {
			BeforeEach(func() {
				timeOfDay = libtime.TimeOfDayFromTime(ParseTime("2023-05-02T13:45:59.123456Z"))
			})
			It("returns correct string", func() {
				Expect(result).To(Equal("13:45:59.123456Z"))
			})
		})
		Context("without nano", func() {
			BeforeEach(func() {
				timeOfDay = libtime.TimeOfDayFromTime(ParseTime("2023-05-02T13:45:59Z"))
			})
			It("returns correct string", func() {
				Expect(result).To(Equal("13:45:59Z"))
			})
		})
	})
	Context("MarshalJSON", func() {
		var bytes []byte
		BeforeEach(func() {
			timeOfDay = libtime.TimeOfDayFromTime(ParseTime("2023-05-02T13:45:59.123456Z"))
		})
		JustBeforeEach(func() {
			bytes, err = timeOfDay.MarshalJSON()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(string(bytes)).To(Equal(`"13:45:59.123456Z"`))
		})
	})
	DescribeTable("UnmarshalJSON",
		func(input string, expected string, expectError bool) {
			timeOfDay = libtime.TimeOfDay{}
			err = timeOfDay.UnmarshalJSON([]byte(input))
			if expectError {
				Expect(err).NotTo(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(timeOfDay.String()).To(Equal(expected))
			}
		},
		Entry("error", `"banana"`, ``, true),
		Entry("with tz", `"13:45:59Z"`, `13:45:59Z`, false),
		Entry("with tz and ns", `"13:45:59.123456Z"`, `13:45:59.123456Z`, false),
		Entry("hour:min tz", `"13:45Z"`, `13:45:00Z`, false),
		Entry("hour:min", `"13:45"`, `13:45:00Z`, false),
		Entry("without tz", `"13:45:59"`, `13:45:59Z`, false),
		Entry("without tz and ns", `"13:45:59.123456"`, `13:45:59.123456Z`, false),
		Entry("datetime with tz", `"2023-10-02T13:45:59Z"`, `13:45:59Z`, false),
		Entry("datetime with tz and ns", `"2023-10-02T13:45:59.123456Z"`, `13:45:59.123456Z`, false),
	)
})

func ParseTime(timeString string) time.Time {
	result, err := time.Parse(time.RFC3339, timeString)
	Expect(err).To(BeNil())
	return result
}
