// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"context"

	. "github.com/onsi/gomega"

	libtime "github.com/bborbe/time"
)

func ParseDuration(value interface{}) libtime.Duration {
	result, err := libtime.ParseDuration(context.Background(), value)
	Expect(err).To(BeNil())
	return *result
}
