// Copyright (c) 2024 Benjamin Borbe All rights reserved.
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

var _ = DescribeTable("ParseDuration",
	func(input string, expectedDuration time.Duration, expectedError bool) {
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
	Entry("ns", "1ns", time.Nanosecond, false),
	Entry("us", "1us", time.Microsecond, false),
	Entry("ms", "1ms", time.Millisecond, false),
	Entry("second", "1s", time.Second, false),
	Entry("minute", "1m", time.Minute, false),
	Entry("hour", "1h", time.Hour, false),
	Entry("day", "1d", 24*time.Hour, false),
	Entry("week", "1w", 7*24*time.Hour, false),
	Entry("combined", "1h30m", 90*time.Minute, false),
	Entry("negative", "-1h30m", -90*time.Minute, false),
	Entry("dot", "1.5h", 90*time.Minute, false),
	Entry("hello", "hello", time.Duration(0), true),
	Entry("hello1d", "hello1d", time.Duration(0), true),
)
