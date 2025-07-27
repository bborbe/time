// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"context"

	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

func ParseDate(value interface{}) libtime.Date {
	result, err := libtime.ParseDate(context.Background(), value)
	Expect(err).To(BeNil())
	return *result
}
