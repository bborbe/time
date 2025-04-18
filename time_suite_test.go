// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"

	libtime "github.com/bborbe/time"
)

//go:generate go run -mod=mod github.com/maxbrunsfeld/counterfeiter/v6 -generate
func TestSuite(t *testing.T) {
	time.Local = time.UTC
	format.TruncatedDiff = false
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

func ParseTimeOfDay(value interface{}) libtime.TimeOfDay {
	result, err := libtime.ParseTimeOfDay(context.Background(), value)
	Expect(err).To(BeNil())
	return *result
}

func ParseTime(input string) time.Time {
	result, err := time.Parse(time.RFC3339, input)
	Expect(err).To(BeNil())
	return result
}

func ParseDateTime(input string) libtime.DateTime {
	result, err := libtime.ParseDateTime(context.Background(), input)
	Expect(err).To(BeNil())
	return *result
}

func ParseUnixTime(input string) libtime.UnixTime {
	result, err := libtime.ParseUnixTime(context.Background(), input)
	Expect(err).To(BeNil())
	return *result
}

func ParseDate(input string) libtime.Date {
	result, err := libtime.ParseDate(context.Background(), input)
	Expect(err).To(BeNil())
	return *result
}
