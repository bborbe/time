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
		var now time.Time
		BeforeEach(func() {
			snapshotTime = libtime.Date{}
			now = time.Unix(1731169783, 0)
			libtime.Now = func() time.Time {
				return now
			}
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
		Context("with value NOW", func() {
			BeforeEach(func() {
				value = `"NOW"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(snapshotTime.Time().Format(time.DateOnly)).To(Equal(`2024-11-09`))
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
				Expect(snapshotTime.Time().Format(time.DateOnly)).To(Equal(`2024-10-26`))
			})
		})
		Context("with value NOW+1d", func() {
			BeforeEach(func() {
				value = `"NOW+1d"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content (now plus 1 day)", func() {
				Expect(snapshotTime.Time().Format(time.DateOnly)).To(Equal(`2024-11-10`))
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
	DescribeTable(
		"ComparePtr",
		func(a *libtime.Date, b *libtime.Date, expectedResult int) {
			Expect(a.ComparePtr(b)).To(Equal(expectedResult))
		},
		Entry(
			"equal",
			libtime.Date(time.Unix(1000, 0)).Ptr(),
			libtime.Date(time.Unix(1000, 0)).Ptr(),
			0,
		),
		Entry(
			"less",
			libtime.Date(time.Unix(999, 0)).Ptr(),
			libtime.Date(time.Unix(1000, 0)).Ptr(),
			-1,
		),
		Entry(
			"greater",
			libtime.Date(time.Unix(1000, 0)).Ptr(),
			libtime.Date(time.Unix(999, 0)).Ptr(),
			1,
		),
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
	Context("TimePtr", func() {
		var dateTime *libtime.Date
		var timePtr *time.Time
		BeforeEach(func() {
			dateTime = libtime.Date(time.Unix(1000, 0)).Ptr()
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
		var dateTime libtime.Date
		var result libtime.Date
		var days int
		var months int
		var years int
		BeforeEach(func() {
			years = 0
			months = 0
			days = 0
			dateTime = ParseDate("2024-12-24")
		})
		JustBeforeEach(func() {
			result = dateTime.AddDate(years, months, days)
		})
		Context("add nothing", func() {
			It("returns the date time", func() {
				Expect(result.String()).To(Equal("2024-12-24"))
			})
		})
		Context("add +1 month", func() {
			BeforeEach(func() {
				months = 1
			})
			It("returns the date time", func() {
				Expect(result.String()).To(Equal("2025-01-24"))
			})
		})
		Context("add -1 month", func() {
			BeforeEach(func() {
				months = -1
			})
			It("returns the date time", func() {
				Expect(result.String()).To(Equal("2024-11-24"))
			})
		})
	})
	Context("IsZero", func() {
		Context("zero time", func() {
			It("returns true", func() {
				var date libtime.Date
				Expect(date.IsZero()).To(BeTrue())
			})
		})
		Context("non-zero time", func() {
			It("returns false", func() {
				date := libtime.ToDate(time.Unix(1687161394, 0))
				Expect(date.IsZero()).To(BeFalse())
			})
		})
	})
	Context("NewDate", func() {
		var result libtime.Date
		BeforeEach(func() {
			result = libtime.NewDate(2023, time.June, 19, 7, 56, 34, 0, time.UTC)
		})
		It("creates correct Date", func() {
			Expect(result.Year()).To(Equal(2023))
			Expect(result.Month()).To(Equal(time.June))
			Expect(result.Day()).To(Equal(19))
		})
		It("matches time.Date behavior", func() {
			expected := time.Date(2023, time.June, 19, 7, 56, 34, 0, time.UTC)
			Expect(result.Time()).To(Equal(expected))
		})
	})
	Context("JSON Regression Tests - Phase 1", func() {
		Context("MarshalJSON *Date nil pointer", func() {
			var nilDate *libtime.Date
			var bytes []byte
			JustBeforeEach(func() {
				bytes, err = json.Marshal(nilDate)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns null", func() {
				Expect(string(bytes)).To(Equal(`null`))
			})
		})
		Context("UnmarshalJSON round-trip with date-only format", func() {
			var originalDate libtime.Date
			var unmarshaledDate libtime.Date
			var jsonBytes []byte
			var remarshaledBytes []byte
			BeforeEach(func() {
				originalDate = libtime.Date(time.Unix(1687161394, 0))
			})
			JustBeforeEach(func() {
				// Marshal original date
				jsonBytes, err = originalDate.MarshalJSON()
				Expect(err).To(BeNil())
				// Unmarshal into new date
				err = unmarshaledDate.UnmarshalJSON(jsonBytes)
				Expect(err).To(BeNil())
				// Re-marshal to verify round-trip
				remarshaledBytes, err = unmarshaledDate.MarshalJSON()
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("unmarshals to correct date", func() {
				Expect(unmarshaledDate.String()).To(Equal("2023-06-19"))
			})
			It("re-marshals to identical JSON", func() {
				Expect(string(remarshaledBytes)).To(Equal(string(jsonBytes)))
				Expect(string(remarshaledBytes)).To(Equal(`"2023-06-19"`))
			})
		})
		Context("UnmarshalJSON with RFC3339 format", func() {
			var date libtime.Date
			BeforeEach(func() {
				date = libtime.Date{}
			})
			JustBeforeEach(func() {
				err = date.UnmarshalJSON([]byte(`"2023-06-19T07:56:34Z"`))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("extracts date-only portion", func() {
				Expect(date.String()).To(Equal("2023-06-19"))
			})
			It("marshals back as date-only", func() {
				bytes, err := date.MarshalJSON()
				Expect(err).To(BeNil())
				Expect(string(bytes)).To(Equal(`"2023-06-19"`))
			})
		})
		Context("JSON struct round-trip", func() {
			type TestStruct struct {
				Date    libtime.Date  `json:"date"`
				DatePtr *libtime.Date `json:"datePtr"`
			}
			var original TestStruct
			var unmarshaled TestStruct
			var jsonBytes []byte
			BeforeEach(func() {
				original = TestStruct{
					Date:    libtime.Date(time.Unix(1687161394, 0)),
					DatePtr: libtime.Date(time.Unix(1687161394, 0)).Ptr(),
				}
			})
			JustBeforeEach(func() {
				// Marshal
				jsonBytes, err = json.Marshal(original)
				Expect(err).To(BeNil())
				// Unmarshal
				err = json.Unmarshal(jsonBytes, &unmarshaled)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("round-trips Date field correctly", func() {
				Expect(unmarshaled.Date.String()).To(Equal("2023-06-19"))
				Expect(unmarshaled.Date.String()).To(Equal(original.Date.String()))
			})
			It("round-trips DatePtr field correctly", func() {
				Expect(unmarshaled.DatePtr).NotTo(BeNil())
				Expect(unmarshaled.DatePtr.String()).To(Equal("2023-06-19"))
				Expect(unmarshaled.DatePtr.String()).To(Equal(original.DatePtr.String()))
			})
			It("marshals to expected JSON format", func() {
				Expect(string(jsonBytes)).To(ContainSubstring(`"date":"2023-06-19"`))
				Expect(string(jsonBytes)).To(ContainSubstring(`"datePtr":"2023-06-19"`))
			})
		})
	})
	Context("MarshalText - Phase 2", func() {
		var date libtime.Date
		var textBytes []byte
		JustBeforeEach(func() {
			textBytes, err = date.MarshalText()
		})
		Context("non-zero date", func() {
			BeforeEach(func() {
				date = libtime.Date(time.Unix(1687161394, 0))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns date-only format without quotes", func() {
				Expect(string(textBytes)).To(Equal("2023-06-19"))
			})
		})
		Context("zero date", func() {
			BeforeEach(func() {
				date = libtime.Date{}
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns nil bytes", func() {
				Expect(textBytes).To(BeNil())
			})
		})
	})
	Context("UnmarshalText - Phase 2", func() {
		var date libtime.Date
		var textInput []byte
		JustBeforeEach(func() {
			err = date.UnmarshalText(textInput)
		})
		Context("date-only format", func() {
			BeforeEach(func() {
				textInput = []byte("2023-06-19")
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("parses correctly", func() {
				Expect(date.String()).To(Equal("2023-06-19"))
			})
		})
		Context("empty string", func() {
			BeforeEach(func() {
				textInput = []byte("")
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("sets to zero date", func() {
				Expect(date.IsZero()).To(BeTrue())
			})
		})
		Context("RFC3339 format", func() {
			BeforeEach(func() {
				textInput = []byte("2023-06-19T07:56:34Z")
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("extracts date portion", func() {
				Expect(date.String()).To(Equal("2023-06-19"))
			})
		})
	})
	Context("YAML round-trip - Phase 2", func() {
		type TestStruct struct {
			Date    libtime.Date  `yaml:"date"`
			DatePtr *libtime.Date `yaml:"datePtr"`
		}
		var original TestStruct
		var unmarshaled TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			original = TestStruct{
				Date:    libtime.Date(time.Unix(1687161394, 0)),
				DatePtr: libtime.Date(time.Unix(1687161394, 0)).Ptr(),
			}
		})
		JustBeforeEach(func() {
			// Marshal to YAML
			yamlBytes, err = yaml.Marshal(original)
			Expect(err).To(BeNil())
			// Unmarshal from YAML
			err = yaml.Unmarshal(yamlBytes, &unmarshaled)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("marshals to date-only format (not timestamp)", func() {
			yamlString := string(yamlBytes)
			// YAML may quote or not quote strings - both are valid
			Expect(yamlString).To(MatchRegexp(`date: "?2023-06-19"?`))
			Expect(yamlString).To(MatchRegexp(`datePtr: "?2023-06-19"?`))
			// Ensure it's NOT a full timestamp
			Expect(yamlString).NotTo(ContainSubstring("T07:56:34"))
		})
		It("round-trips Date field correctly", func() {
			Expect(unmarshaled.Date.String()).To(Equal("2023-06-19"))
			Expect(unmarshaled.Date.String()).To(Equal(original.Date.String()))
		})
		It("round-trips DatePtr field correctly", func() {
			Expect(unmarshaled.DatePtr).NotTo(BeNil())
			Expect(unmarshaled.DatePtr.String()).To(Equal("2023-06-19"))
			Expect(unmarshaled.DatePtr.String()).To(Equal(original.DatePtr.String()))
		})
	})
	Context("YAML omitempty - Phase 2", func() {
		type TestStruct struct {
			Date        libtime.Date  `yaml:"date,omitempty"`
			DatePtr     *libtime.Date `yaml:"datePtr,omitempty"`
			DateNonZero libtime.Date  `yaml:"dateNonZero,omitempty"`
		}
		var testStruct TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			testStruct = TestStruct{
				Date:        libtime.Date{},                         // zero
				DatePtr:     nil,                                    // nil
				DateNonZero: libtime.Date(time.Unix(1687161394, 0)), // non-zero
			}
		})
		JustBeforeEach(func() {
			yamlBytes, err = yaml.Marshal(testStruct)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("omits zero Date field with omitempty", func() {
			yamlString := string(yamlBytes)
			// Zero date with omitempty should be omitted (YAML treats zero time specially)
			// Note: This behavior depends on YAML implementation
			Expect(yamlString).NotTo(ContainSubstring("date:"))
		})
		It("omits nil *Date field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("datePtr:"))
		})
		It("includes non-zero Date field", func() {
			yamlString := string(yamlBytes)
			// YAML may quote or not quote strings - both are valid
			Expect(yamlString).To(MatchRegexp(`dateNonZero: "?2023-06-19"?`))
		})
	})
	Context("struct marshal regression - Phase 1", func() {
		Context("A. JSON Marshal — non-zero values set for Field and FieldPtr only", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.Date  `json:"field"`
					FieldPtr     *libtime.Date `json:"fieldPtr"`
					FieldOmit    libtime.Date  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Date `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    libtime.Date(time.Unix(1687161394, 0)),
					FieldPtr: libtime.Date(time.Unix(1687161394, 0)).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns exact JSON output", func() {
				Expect(
					string(jsonBytes),
				).To(Equal(`{"field":"2023-06-19","fieldPtr":"2023-06-19","fieldOmit":null}`))
			})
		})
		Context("B. JSON Unmarshal — round-trip", func() {
			var original struct {
				Field        libtime.Date  `json:"field"`
				FieldPtr     *libtime.Date `json:"fieldPtr"`
				FieldOmit    libtime.Date  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.Date `json:"fieldPtrOmit,omitempty"`
			}
			var unmarshaled struct {
				Field        libtime.Date  `json:"field"`
				FieldPtr     *libtime.Date `json:"fieldPtr"`
				FieldOmit    libtime.Date  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.Date `json:"fieldPtrOmit,omitempty"`
			}
			var jsonBytes []byte
			BeforeEach(func() {
				original = struct {
					Field        libtime.Date  `json:"field"`
					FieldPtr     *libtime.Date `json:"fieldPtr"`
					FieldOmit    libtime.Date  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Date `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    libtime.Date(time.Unix(1687161394, 0)),
					FieldPtr: libtime.Date(time.Unix(1687161394, 0)).Ptr(),
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
					Field        libtime.Date  `json:"field"`
					FieldPtr     *libtime.Date `json:"fieldPtr"`
					FieldOmit    libtime.Date  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Date `json:"fieldPtrOmit,omitempty"`
				}{
					Field:        libtime.Date(time.Unix(1687161394, 0)),
					FieldPtr:     libtime.Date(time.Unix(1687161394, 0)).Ptr(),
					FieldOmit:    libtime.Date(time.Unix(1687161394, 0)),
					FieldPtrOmit: libtime.Date(time.Unix(1687161394, 0)).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("all fields appear in output", func() {
				jsonStr := string(jsonBytes)
				Expect(jsonStr).To(ContainSubstring(`"field":"2023-06-19"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtr":"2023-06-19"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldOmit":"2023-06-19"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtrOmit":"2023-06-19"`))
			})
		})
		Context("D. JSON Marshal — all fields zero/nil", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.Date  `json:"field"`
					FieldPtr     *libtime.Date `json:"fieldPtr"`
					FieldOmit    libtime.Date  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Date `json:"fieldPtrOmit,omitempty"`
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
