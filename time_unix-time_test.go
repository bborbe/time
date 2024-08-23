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

var _ = Describe("UnixTime", func() {
	var err error
	var snapshotTime libtime.UnixTime
	Context("MarshalJSON", func() {
		var bytes []byte
		BeforeEach(func() {
			snapshotTime = libtime.UnixTime(time.Unix(1687161394, 0))
		})
		JustBeforeEach(func() {
			bytes, err = snapshotTime.MarshalJSON()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(string(bytes)).To(Equal(`1687161394`))
		})
	})
	Context("UnmarshalJSON", func() {
		BeforeEach(func() {
			snapshotTime = libtime.UnixTime{}
		})
		JustBeforeEach(func() {
			err = snapshotTime.UnmarshalJSON([]byte(`1687161394`))
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(snapshotTime.Time().Format(time.RFC3339Nano)).To(Equal(`2023-06-19T07:56:34Z`))
		})
	})
})
