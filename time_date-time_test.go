// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	libtime "github.com/bborbe/time"
)

var _ = Describe("DateTime", func() {
	var err error
	var ctx context.Context
	var now time.Time
	BeforeEach(func() {
		ctx = context.Background()
		now = time.Unix(1731169783, 0)
		libtime.Now = func() time.Time {
			return now
		}
	})
	Context("MarshalBinary & DateTimeFromBinary", func() {
		var dateTime libtime.DateTime
		var binary []byte
		BeforeEach(func() {
			dateTime = libtime.DateTime(time.Unix(1687161394, 0))
		})
		JustBeforeEach(func() {
			binary, err = dateTime.MarshalBinary()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns binary", func() {
			Expect(binary).NotTo(BeNil())
		})
		It("returns binary", func() {
			dateTimeFromBinary, err := libtime.DateTimeFromBinary(ctx, binary)
			Expect(err).To(BeNil())
			Expect(dateTimeFromBinary).NotTo(BeNil())
			Expect(dateTimeFromBinary.Unix()).To(Equal(int64(1687161394)))
		})
	})
	Context("MarshalJSON", func() {
		var dateTime libtime.DateTime
		var bytes []byte
		JustBeforeEach(func() {
			bytes, err = dateTime.MarshalJSON()
		})
		Context("defined", func() {
			BeforeEach(func() {
				dateTime = libtime.DateTime(time.Unix(1687161394, 0))
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
				dateTime = libtime.DateTime{}
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
				Expect(
					snapshotTime.Time().Format(time.RFC3339Nano),
				).To(Equal(`2023-06-19T07:56:34Z`))
			})
		})
		Context("with value now", func() {
			BeforeEach(func() {
				value = `"NOW"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(
					snapshotTime.Time().Format(time.RFC3339Nano),
				).To(Equal(`2024-11-09T16:29:43Z`))
			})
		})
		Context("with value NOW-14d", func() {
			BeforeEach(func() {
				value = `"NOW-14d"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content (now minus 14 days)", func() {
				Expect(
					snapshotTime.Time().Format(time.RFC3339Nano),
				).To(Equal(`2024-10-26T16:29:43Z`))
			})
		})
		Context("with value NOW+1h", func() {
			BeforeEach(func() {
				value = `"NOW+1h"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content (now plus 1 hour)", func() {
				Expect(
					snapshotTime.Time().Format(time.RFC3339Nano),
				).To(Equal(`2024-11-09T17:29:43Z`))
			})
		})
		Context("with value NOW-7d", func() {
			BeforeEach(func() {
				value = `"NOW-7d"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content (now minus 7 days)", func() {
				Expect(
					snapshotTime.Time().Format(time.RFC3339Nano),
				).To(Equal(`2024-11-02T16:29:43Z`))
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
	DescribeTable(
		"ParseDateTime",
		func(value interface{}, formatedDate string, expectError bool) {
			result, err := libtime.ParseDateTime(ctx, value)
			if expectError {
				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Format(time.RFC3339Nano)).To(Equal(formatedDate))
			}
		},
		Entry("2023-06-19T07:56:34Z", "2023-06-19T07:56:34Z", "2023-06-19T07:56:34Z", false),
		Entry(
			"2023-06-19T07:56:34.234Z",
			"2023-06-19T07:56:34.234Z",
			"2023-06-19T07:56:34.234Z",
			false,
		),
		Entry("2023-06-19T07:56Z", "2023-06-19T07:56Z", "2023-06-19T07:56:00Z", false),
		Entry("NOW", "NOW", "2024-11-09T16:29:43Z", false),
	)
	DescribeTable(
		"ComparePtr",
		func(a *libtime.DateTime, b *libtime.DateTime, expectedResult int) {
			Expect(a.ComparePtr(b)).To(Equal(expectedResult))
		},
		Entry(
			"equal",
			libtime.DateTime(time.Unix(1000, 0)).Ptr(),
			libtime.DateTime(time.Unix(1000, 0)).Ptr(),
			0,
		),
		Entry(
			"less",
			libtime.DateTime(time.Unix(999, 0)).Ptr(),
			libtime.DateTime(time.Unix(1000, 0)).Ptr(),
			-1,
		),
		Entry(
			"greater",
			libtime.DateTime(time.Unix(1000, 0)).Ptr(),
			libtime.DateTime(time.Unix(999, 0)).Ptr(),
			1,
		),
		Entry("equal", nil, nil, 0),
		Entry("less", nil, libtime.DateTime(time.Unix(1000, 0)).Ptr(), -1),
		Entry("greater", libtime.DateTime(time.Unix(1000, 0)).Ptr(), nil, 1),
	)
	DescribeTable(
		"Compare",
		func(a libtime.DateTime, b libtime.DateTime, expectedResult int) {
			Expect(a.Compare(b)).To(Equal(expectedResult))
		},
		Entry(
			"equal",
			libtime.DateTime(time.Unix(1000, 0)),
			libtime.DateTime(time.Unix(1000, 0)),
			0,
		),
		Entry(
			"less",
			libtime.DateTime(time.Unix(999, 0)),
			libtime.DateTime(time.Unix(1000, 0)),
			-1,
		),
		Entry(
			"greater",
			libtime.DateTime(time.Unix(1000, 0)),
			libtime.DateTime(time.Unix(999, 0)),
			1,
		),
	)
	Context("TimePtr", func() {
		var dateTime *libtime.DateTime
		var timePtr *time.Time
		BeforeEach(func() {
			dateTime = libtime.DateTime(time.Unix(1000, 0)).Ptr()
		})
		JustBeforeEach(func() {
			timePtr = dateTime.TimePtr()
		})
		Context("datetime not nil", func() {
			It("returns timePtr", func() {
				Expect(timePtr).NotTo(BeNil())
			})
		})
		Context("datetime nil", func() {
			BeforeEach(func() {
				dateTime = nil
			})
			It("returns not timePtr", func() {
				Expect(timePtr).To(BeNil())
			})
		})
	})
	Context("AddDate", func() {
		var dateTime libtime.DateTime
		var result libtime.DateTime
		var days int
		var months int
		var years int
		BeforeEach(func() {
			years = 0
			months = 0
			days = 0
			dateTime = ParseDateTime("2024-12-24T20:15:59Z")
		})
		JustBeforeEach(func() {
			result = dateTime.AddDate(years, months, days)
		})
		Context("add nothing", func() {
			It("returns the date time", func() {
				Expect(result.String()).To(Equal("2024-12-24T20:15:59Z"))
			})
		})
		Context("add +1 month", func() {
			BeforeEach(func() {
				months = 1
			})
			It("returns the date time", func() {
				Expect(result.String()).To(Equal("2025-01-24T20:15:59Z"))
			})
		})
		Context("add -1 month", func() {
			BeforeEach(func() {
				months = -1
			})
			It("returns the date time", func() {
				Expect(result.String()).To(Equal("2024-11-24T20:15:59Z"))
			})
		})
	})
	Context("IsZero", func() {
		Context("zero time", func() {
			It("returns true", func() {
				var dateTime libtime.DateTime
				Expect(dateTime.IsZero()).To(BeTrue())
			})
		})
		Context("non-zero time", func() {
			It("returns false", func() {
				dateTime := libtime.DateTime(time.Unix(1687161394, 0))
				Expect(dateTime.IsZero()).To(BeFalse())
			})
		})
	})
	Context("NewDateTime", func() {
		var result libtime.DateTime
		BeforeEach(func() {
			result = libtime.NewDateTime(2023, time.June, 19, 7, 56, 34, 0, time.UTC)
		})
		It("creates correct DateTime", func() {
			Expect(result.Year()).To(Equal(2023))
			Expect(result.Month()).To(Equal(time.June))
			Expect(result.Day()).To(Equal(19))
			Expect(result.Hour()).To(Equal(7))
			Expect(result.Minute()).To(Equal(56))
			Expect(result.Second()).To(Equal(34))
			Expect(result.Nanosecond()).To(Equal(0))
			Expect(result.Time().Location()).To(Equal(time.UTC))
		})
		It("matches time.Date behavior", func() {
			expected := time.Date(2023, time.June, 19, 7, 56, 34, 0, time.UTC)
			Expect(result.Time()).To(Equal(expected))
		})
	})
	Context("YAML round-trip - Phase 2", func() {
		type TestStruct struct {
			DateTime    libtime.DateTime  `yaml:"dateTime"`
			DateTimePtr *libtime.DateTime `yaml:"dateTimePtr"`
		}
		var original TestStruct
		var unmarshaled TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			original = TestStruct{
				DateTime:    libtime.DateTime(time.Unix(1687161394, 0)),
				DateTimePtr: libtime.DateTime(time.Unix(1687161394, 0)).Ptr(),
			}
		})
		JustBeforeEach(func() {
			// Marshal to YAML
			var err error
			yamlBytes, err = yaml.Marshal(original)
			Expect(err).To(BeNil())
			// Unmarshal from YAML
			err = yaml.Unmarshal(yamlBytes, &unmarshaled)
			Expect(err).To(BeNil())
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("marshals to RFC3339Nano format (not integer)", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`dateTime: "?2023-06-19T07:56:34Z"?`))
			Expect(yamlString).To(MatchRegexp(`dateTimePtr: "?2023-06-19T07:56:34Z"?`))
		})
		It("round-trips DateTime field correctly", func() {
			Expect(unmarshaled.DateTime.String()).To(Equal("2023-06-19T07:56:34Z"))
			Expect(unmarshaled.DateTime.String()).To(Equal(original.DateTime.String()))
		})
		It("round-trips DateTimePtr field correctly", func() {
			Expect(unmarshaled.DateTimePtr).NotTo(BeNil())
			Expect(unmarshaled.DateTimePtr.String()).To(Equal("2023-06-19T07:56:34Z"))
			Expect(unmarshaled.DateTimePtr.String()).To(Equal(original.DateTimePtr.String()))
		})
	})
	Context("YAML omitempty - Phase 2", func() {
		type TestStruct struct {
			DateTime        libtime.DateTime  `yaml:"dateTime,omitempty"`
			DateTimePtr     *libtime.DateTime `yaml:"dateTimePtr,omitempty"`
			DateTimeNonZero libtime.DateTime  `yaml:"dateTimeNonZero,omitempty"`
		}
		var testStruct TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			testStruct = TestStruct{
				DateTime:        libtime.DateTime{},                         // zero
				DateTimePtr:     nil,                                        // nil
				DateTimeNonZero: libtime.DateTime(time.Unix(1687161394, 0)), // non-zero
			}
		})
		JustBeforeEach(func() {
			yamlBytes, err = yaml.Marshal(testStruct)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("omits zero DateTime field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("dateTime:"))
		})
		It("omits nil *DateTime field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("dateTimePtr:"))
		})
		It("includes non-zero DateTime field", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`dateTimeNonZero: "?2023-06-19T07:56:34Z"?`))
		})
	})
	Context("struct marshal regression - Phase 1", func() {
		Context("A. JSON Marshal — non-zero values set for Field and FieldPtr only", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.DateTime  `json:"field"`
					FieldPtr     *libtime.DateTime `json:"fieldPtr"`
					FieldOmit    libtime.DateTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.DateTime `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    libtime.DateTime(time.Unix(1687161394, 0)),
					FieldPtr: libtime.DateTime(time.Unix(1687161394, 0)).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns exact JSON output", func() {
				Expect(
					string(jsonBytes),
				).To(Equal(`{"field":"2023-06-19T07:56:34Z","fieldPtr":"2023-06-19T07:56:34Z","fieldOmit":null}`))
			})
		})
		Context("B. JSON Unmarshal — round-trip", func() {
			var original struct {
				Field        libtime.DateTime  `json:"field"`
				FieldPtr     *libtime.DateTime `json:"fieldPtr"`
				FieldOmit    libtime.DateTime  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.DateTime `json:"fieldPtrOmit,omitempty"`
			}
			var unmarshaled struct {
				Field        libtime.DateTime  `json:"field"`
				FieldPtr     *libtime.DateTime `json:"fieldPtr"`
				FieldOmit    libtime.DateTime  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.DateTime `json:"fieldPtrOmit,omitempty"`
			}
			var jsonBytes []byte
			BeforeEach(func() {
				original = struct {
					Field        libtime.DateTime  `json:"field"`
					FieldPtr     *libtime.DateTime `json:"fieldPtr"`
					FieldOmit    libtime.DateTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.DateTime `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    libtime.DateTime(time.Unix(1687161394, 0)),
					FieldPtr: libtime.DateTime(time.Unix(1687161394, 0)).Ptr(),
				}
			})
			JustBeforeEach(func() {
				jsonBytes, err = json.Marshal(original)
				Expect(err).To(BeNil())
				err = json.Unmarshal(jsonBytes, &unmarshaled)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("round-trips Field correctly", func() {
				Expect(unmarshaled.Field.String()).To(Equal(original.Field.String()))
			})
			It("round-trips FieldPtr correctly", func() {
				Expect(unmarshaled.FieldPtr).NotTo(BeNil())
				Expect(unmarshaled.FieldPtr.String()).To(Equal(original.FieldPtr.String()))
			})
			It("zero fields remain zero", func() {
				Expect(unmarshaled.FieldOmit.IsZero()).To(BeTrue())
				Expect(unmarshaled.FieldPtrOmit).To(BeNil())
			})
		})
		Context("C. JSON Marshal — all fields set (non-zero)", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.DateTime  `json:"field"`
					FieldPtr     *libtime.DateTime `json:"fieldPtr"`
					FieldOmit    libtime.DateTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.DateTime `json:"fieldPtrOmit,omitempty"`
				}{
					Field:        libtime.DateTime(time.Unix(1687161394, 0)),
					FieldPtr:     libtime.DateTime(time.Unix(1687161394, 0)).Ptr(),
					FieldOmit:    libtime.DateTime(time.Unix(1687161394, 0)),
					FieldPtrOmit: libtime.DateTime(time.Unix(1687161394, 0)).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("all fields appear in output", func() {
				jsonStr := string(jsonBytes)
				Expect(jsonStr).To(ContainSubstring(`"field":"2023-06-19T07:56:34Z"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtr":"2023-06-19T07:56:34Z"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldOmit":"2023-06-19T07:56:34Z"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtrOmit":"2023-06-19T07:56:34Z"`))
			})
		})
		Context("D. JSON Marshal — all fields zero/nil", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.DateTime  `json:"field"`
					FieldPtr     *libtime.DateTime `json:"fieldPtr"`
					FieldOmit    libtime.DateTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.DateTime `json:"fieldPtrOmit,omitempty"`
				}{}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("non-omitempty fields produce null, omitempty fields produce null", func() {
				Expect(
					string(jsonBytes),
				).To(Equal(`{"field":null,"fieldPtr":null,"fieldOmit":null}`))
			})
		})
	})
})
