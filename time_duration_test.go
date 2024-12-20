// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"

	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("ParseDuration",
	func(input string, expectedDuration libtime.Duration, expectedError bool) {
		duration, err := libtime.ParseDuration(context.Background(), input)
		if expectedError {
			Expect(err).NotTo(BeNil())
			Expect(duration).To(BeNil())
		} else {
			Expect(err).To(BeNil())
			Expect(duration).NotTo(BeNil())
			Expect(*duration).To(Equal(expectedDuration))
		}
	},
	Entry("without unit", "1337", libtime.Duration(1337), false),
	Entry("ns", "1ns", libtime.Nanosecond, false),
	Entry("us", "1us", libtime.Microsecond, false),
	Entry("ms", "1ms", libtime.Millisecond, false),
	Entry("second", "1s", libtime.Second, false),
	Entry("minute", "1m", libtime.Minute, false),
	Entry("hour", "1h", libtime.Hour, false),
	Entry("day", "1d", 24*libtime.Hour, false),
	Entry("week", "1w", 7*24*libtime.Hour, false),
	Entry("combined", "1h30m", 90*libtime.Minute, false),
	Entry("negative", "-1h30m", -90*libtime.Minute, false),
	Entry("dot", "1.5h", 90*libtime.Minute, false),
	Entry("hello", "hello", libtime.Duration(0), true),
	Entry("hello1d", "hello1d", libtime.Duration(0), true),
)

var _ = Describe("Duration", func() {
	var _ = DescribeTable("String",
		func(inputDuration libtime.Duration, expectedOutput string) {
			Expect(inputDuration.String()).To(Equal(expectedOutput))
		},
		Entry("30s", 30*libtime.Second, "30s"),
		Entry("59m30s", 59*libtime.Minute+30*libtime.Second, "59m30s"),
		Entry("23h59m30s", 23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "23h59m30s"),
		Entry("5d23h59m30s", 5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "5d23h59m30s"),
		Entry("10w5d23h59m30s", 10*libtime.Week+5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "10w5d23h59m30s"),
	)

	var _ = DescribeTable("MarshalJSON",
		func(inputDuration libtime.Duration, expectedOutput string, expectError bool) {
			bytes, err := inputDuration.MarshalJSON()
			if expectError {
				Expect(err).NotTo(BeNil())
				Expect(bytes).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(string(bytes)).To(Equal(expectedOutput))
			}
		},
		Entry("0", libtime.Duration(0), `"0s"`, false),
		Entry("30s", 30*libtime.Second, `"30s"`, false),
		Entry("59m30s", 59*libtime.Minute+30*libtime.Second, `"59m30s"`, false),
		Entry("23h59m30s", 23*libtime.Hour+59*libtime.Minute+30*libtime.Second, `"23h59m30s"`, false),
		Entry("143h59m30s", 5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second, `"143h59m30s"`, false),
		Entry("1823h59m30s", 10*libtime.Week+5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second, `"1823h59m30s"`, false),
	)

	var _ = DescribeTable("String",
		func(inputDuration libtime.Duration, expectedOutput string) {
			Expect(inputDuration.String()).To(Equal(expectedOutput))
		},
		Entry("1w", libtime.Week, "1w"),
		Entry("1d", libtime.Day, "1d"),
		Entry("1h", libtime.Hour, "1h"),
		Entry("1m", libtime.Minute, "1m"),
		Entry("1s", libtime.Second, "1s"),
		Entry("1ms", libtime.Millisecond, "1ms"),
		Entry("1µs", libtime.Microsecond, "1µs"),
		Entry("1ns", libtime.Nanosecond, "1ns"),
		Entry("0", libtime.Duration(0), "0s"),
		Entry("1w1ns", libtime.Week+libtime.Nanosecond, "1w1ns"),
		Entry("59m30s", 59*libtime.Minute+30*libtime.Second, "59m30s"),
		Entry("23h59m30s", 23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "23h59m30s"),
		Entry("5d23h59m30s", 5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "5d23h59m30s"),
		Entry("10w5d23h59m30s", 10*libtime.Week+5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "10w5d23h59m30s"),
	)

	Context("UnmarshalJSON", func() {
		var err error
		var duration libtime.Duration
		var value string
		BeforeEach(func() {
			duration = 0
		})
		JustBeforeEach(func() {
			err = duration.UnmarshalJSON([]byte(value))
		})
		Context("with string value", func() {
			BeforeEach(func() {
				value = `"1337"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(libtime.Duration(1337)))
			})
		})
		Context("with number value", func() {
			BeforeEach(func() {
				value = `1337`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(libtime.Duration(1337)))
			})
		})
		Context("with duration value", func() {
			BeforeEach(func() {
				value = `"1h"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(libtime.Hour))
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
				Expect(duration).To(Equal(libtime.Duration(0)))
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
				Expect(duration).To(Equal(libtime.Duration(0)))
			})
		})
	})
})
