// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

var _ = DescribeTable("Compare",
	func(a time.Time, b time.Time, expectedResult int) {
		Expect(libtime.Compare(a, b)).To(Equal(expectedResult))
	},
	Entry("equal", time.Unix(1000, 0), time.Unix(1000, 0), 0),
	Entry("less", time.Unix(999, 0), time.Unix(1000, 0), -1),
	Entry("greater", time.Unix(1000, 0), time.Unix(999, 0), 1),
)
