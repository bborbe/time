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

var _ = Describe("Location", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	DescribeTable("LoadLocation",
		func(input string, expectedLocation string, expectedError bool) {
			location, err := libtime.LoadLocation(ctx, input)
			if expectedError {
				Expect(err).NotTo(BeNil())
				Expect(location).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(location).NotTo(BeNil())
				Expect(location.String()).To(Equal(expectedLocation))
			}
		},
		Entry("UTC", "UTC", "UTC", false),
		Entry("Europe/Berlin", "Europe/Berlin", "Europe/Berlin", false),
		Entry("Banana", "Banana", "", true),
	)
	DescribeTable("ParseLocation",
		func(input any, expectedLocation string, expectedError bool) {
			location, err := libtime.ParseLocation(ctx, input)
			if expectedError {
				Expect(err).NotTo(BeNil())
				Expect(location).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(location).NotTo(BeNil())
				Expect(location.String()).To(Equal(expectedLocation))
			}
		},
		Entry("UTC", "UTC", "UTC", false),
		Entry("Europe/Berlin", "Europe/Berlin", "Europe/Berlin", false),
		Entry("Banana", "Banana", "", true),
	)
})
