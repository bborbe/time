// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"time"

	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TimeOfDay", func() {
	var err error
	var timeOfDay libtime.TimeOfDay
	BeforeEach(func() {
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
