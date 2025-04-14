// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

var _ = Describe("Weekday", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	DescribeTable("Validate",
		func(input libtime.Weekday, expectedError bool) {
			if expectedError {
				Expect(input.Validate(ctx)).NotTo(BeNil())
			} else {
				Expect(input.Validate(ctx)).To(BeNil())
			}
		},
		Entry("Monday", libtime.Monday, false),
		Entry("Tuesday", libtime.Tuesday, false),
		Entry("Wednesday", libtime.Wednesday, false),
		Entry("Thursday", libtime.Thursday, false),
		Entry("Friday", libtime.Friday, false),
		Entry("Saturday", libtime.Saturday, false),
		Entry("Sunday", libtime.Sunday, false),
		Entry("invalid", libtime.Weekday(1337), true),
	)
})
