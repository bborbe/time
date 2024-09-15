// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DateTime", func() {
	var err error
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("MarshalJSON", func() {
		var snapshotTime libtime.DateTime
		var bytes []byte
		JustBeforeEach(func() {
			bytes, err = snapshotTime.MarshalJSON()
		})
		Context("defined", func() {
			BeforeEach(func() {
				snapshotTime = libtime.DateTime(time.Unix(1687161394, 0))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(string(bytes)).To(Equal(`"2023-06-19T07:56:34Z"`))
			})
		})
		Context("undefined", func() {
			BeforeEach(func() {
				snapshotTime = libtime.DateTime{}
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(string(bytes)).To(Equal(`null`))
			})
		})
	})
	Context("json marshal", func() {
		var content string
		JustBeforeEach(func() {
			buf := &bytes.Buffer{}
			encoder := json.NewEncoder(buf)
			encoder.SetIndent("", "  ")

			err = encoder.Encode(struct {
				DateEmpty        libtime.DateTime  `json:"dateEmpty"`
				DatePtrEmpty     *libtime.DateTime `json:"datePtrEmpty"`
				DateOmitEmpty    libtime.DateTime  `json:"dateOmitEmpty,omitempty"`
				DatePtrOmitEmpty *libtime.DateTime `json:"datePtrOmitEmpty,omitempty"`
				Date             libtime.DateTime  `json:"date"`
				DatePtr          *libtime.DateTime `json:"datePtr"`
			}{
				Date:    libtime.DateTime(time.Unix(1687161394, 0)),
				DatePtr: libtime.DateTime(time.Unix(1687161394, 0)).Ptr(),
			})
			content = buf.String()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(content).To(Equal(`{
  "dateEmpty": null,
  "datePtrEmpty": null,
  "dateOmitEmpty": null,
  "date": "2023-06-19T07:56:34Z",
  "datePtr": "2023-06-19T07:56:34Z"
}
`))
		})
	})
	Context("UnmarshalJSON", func() {
		var snapshotTime libtime.DateTime
		var value string
		BeforeEach(func() {
			snapshotTime = libtime.DateTime{}
		})
		JustBeforeEach(func() {
			err = snapshotTime.UnmarshalJSON([]byte(value))
		})
		Context("with value", func() {
			BeforeEach(func() {
				value = `"2023-06-19T07:56:34Z"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(snapshotTime.Time().Format(time.RFC3339Nano)).To(Equal(`2023-06-19T07:56:34Z`))
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
				Expect(snapshotTime.Time().IsZero()).To(BeTrue())
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
				Expect(snapshotTime.Time().IsZero()).To(BeTrue())
			})
		})
	})
	Context("ParseDateTime", func() {
		var value interface{}
		var stdTime *libtime.DateTime
		JustBeforeEach(func() {
			stdTime, err = libtime.ParseDateTime(ctx, value)
		})
		Context("Success", func() {
			BeforeEach(func() {
				value = "2023-06-19T07:56:34Z"
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct time", func() {
				Expect(stdTime).NotTo(BeNil())
				Expect(stdTime.Format(time.RFC3339)).To(Equal("2023-06-19T07:56:34Z"))
			})
		})
	})
	DescribeTable("ComparePtr",
		func(a *libtime.DateTime, b *libtime.DateTime, expectedResult int) {
			Expect(a.ComparePtr(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.DateTime(time.Unix(1000, 0)).Ptr(), libtime.DateTime(time.Unix(1000, 0)).Ptr(), 0),
		Entry("less", libtime.DateTime(time.Unix(999, 0)).Ptr(), libtime.DateTime(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.DateTime(time.Unix(1000, 0)).Ptr(), libtime.DateTime(time.Unix(999, 0)).Ptr(), 1),
		Entry("equal", nil, nil, 0),
		Entry("less", nil, libtime.DateTime(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.DateTime(time.Unix(1000, 0)).Ptr(), nil, 1),
	)
	DescribeTable("Compare",
		func(a libtime.DateTime, b libtime.DateTime, expectedResult int) {
			Expect(a.Compare(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.DateTime(time.Unix(1000, 0)), libtime.DateTime(time.Unix(1000, 0)), 0),
		Entry("less", libtime.DateTime(time.Unix(999, 0)), libtime.DateTime(time.Unix(1000, 0)), -1),
		Entry("greater", libtime.DateTime(time.Unix(1000, 0)), libtime.DateTime(time.Unix(999, 0)), 1),
	)
})