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
	func(input string, expectedDuration time.Duration, expectedError error) {
		duration, err := libtime.ParseDuration(context.Background(), input)
		if expectedError != nil {
			Expect(err).To(Equal(expectedError))
			Expect(duration).To(BeNil())
		} else {
			Expect(err).To(BeNil())
			Expect(duration).NotTo(BeNil())
			Expect(*duration).To(Equal(expectedDuration))
		}
	},
	Entry("minute", "1m", time.Minute, nil),
	Entry("hour", "1h", time.Hour, nil),
	Entry("day", "1d", 24*time.Hour, nil),
	Entry("week", "1w", 7*24*time.Hour, nil),
	Entry("combined", "1h30m", 90*time.Minute, nil),
	Entry("negative", "-1h30m", -90*time.Minute, nil),
	Entry("dot", "1.5h", 90*time.Minute, nil),
)
