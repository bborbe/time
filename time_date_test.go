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

var _ = Describe("Date", func() {
	var err error
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("MarshalBinary & DateFromBinary", func() {
		var date libtime.Date
		var binary []byte
		BeforeEach(func() {
			date = libtime.Date(time.Unix(1687161394, 0))
		})
		JustBeforeEach(func() {
			binary, err = date.MarshalBinary()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns binary", func() {
			Expect(binary).NotTo(BeNil())
		})
		It("returns binary", func() {
			dateFromBinary, err := libtime.DateFromBinary(ctx, binary)
			Expect(err).To(BeNil())
			Expect(dateFromBinary).NotTo(BeNil())
			Expect(dateFromBinary.Unix()).To(Equal(int64(1687161394)))
		})
	})
	Context("MarshalJSON", func() {
		var snapshotTime libtime.Date
		var bytes []byte
		JustBeforeEach(func() {
			bytes, err = snapshotTime.MarshalJSON()
		})
		Context("defined", func() {
			BeforeEach(func() {
				snapshotTime = libtime.Date(time.Unix(1687161394, 0))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(string(bytes)).To(Equal(`"2023-06-19"`))
			})
		})
		Context("undefined", func() {
			BeforeEach(func() {
				snapshotTime = libtime.Date{}
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
				DateEmpty        libtime.Date  `json:"dateEmpty"`
				DatePtrEmpty     *libtime.Date `json:"datePtrEmpty"`
				DateOmitEmpty    libtime.Date  `json:"dateOmitEmpty,omitempty"`
				DatePtrOmitEmpty *libtime.Date `json:"datePtrOmitEmpty,omitempty"`
				Date             libtime.Date  `json:"date"`
				DatePtr          *libtime.Date `json:"datePtr"`
			}{
				Date:    libtime.Date(time.Unix(1687161394, 0)),
				DatePtr: libtime.Date(time.Unix(1687161394, 0)).Ptr(),
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
  "date": "2023-06-19",
  "datePtr": "2023-06-19"
}
`))
		})
	})
	Context("UnmarshalJSON", func() {
		var snapshotTime libtime.Date
		var value string
		BeforeEach(func() {
			snapshotTime = libtime.Date{}
		})
		JustBeforeEach(func() {
			err = snapshotTime.UnmarshalJSON([]byte(value))
		})
		Context("with value", func() {
			BeforeEach(func() {
				value = `"2023-06-19"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(snapshotTime.Time().Format(time.DateOnly)).To(Equal(`2023-06-19`))
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
	Context("ParseDate", func() {
		var value interface{}
		var stdTime *libtime.Date
		JustBeforeEach(func() {
			stdTime, err = libtime.ParseDate(ctx, value)
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
				Expect(stdTime.Format(time.DateOnly)).To(Equal("2023-06-19"))
			})
		})
	})
	DescribeTable("ComparePtr",
		func(a *libtime.Date, b *libtime.Date, expectedResult int) {
			Expect(a.ComparePtr(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.Date(time.Unix(1000, 0)).Ptr(), libtime.Date(time.Unix(1000, 0)).Ptr(), 0),
		Entry("less", libtime.Date(time.Unix(999, 0)).Ptr(), libtime.Date(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.Date(time.Unix(1000, 0)).Ptr(), libtime.Date(time.Unix(999, 0)).Ptr(), 1),
		Entry("equal", nil, nil, 0),
		Entry("less", nil, libtime.Date(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.Date(time.Unix(1000, 0)).Ptr(), nil, 1),
	)
	DescribeTable("Compare",
		func(a libtime.Date, b libtime.Date, expectedResult int) {
			Expect(a.Compare(b)).To(Equal(expectedResult))
		},
		Entry("equal", libtime.Date(time.Unix(1000, 0)), libtime.Date(time.Unix(1000, 0)), 0),
		Entry("less", libtime.Date(time.Unix(999, 0)), libtime.Date(time.Unix(1000, 0)), -1),
		Entry("greater", libtime.Date(time.Unix(1000, 0)), libtime.Date(time.Unix(999, 0)), 1),
	)
})
