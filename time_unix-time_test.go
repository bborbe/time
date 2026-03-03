// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"encoding/json"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	libtime "github.com/bborbe/time"
)

var _ = Describe("UnixTime", func() {
	var err error
	var snapshotTime libtime.UnixTime
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	Context("MarshalBinary & UnixTimeFromBinary", func() {
		var unixTime libtime.UnixTime
		var binary []byte
		BeforeEach(func() {
			unixTime = libtime.UnixTime(time.Unix(1687161394, 0))
		})
		JustBeforeEach(func() {
			binary, err = unixTime.MarshalBinary()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns binary", func() {
			Expect(binary).NotTo(BeNil())
		})
		It("returns binary", func() {
			unixTimeFromBinary, err := libtime.UnixTimeFromBinary(ctx, binary)
			Expect(err).To(BeNil())
			Expect(unixTimeFromBinary).NotTo(BeNil())
			Expect(unixTimeFromBinary.Unix()).To(Equal(int64(1687161394)))
		})
	})
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
	DescribeTable("ParseUnixTime",
		func(input any, expectedDateString string, expectedError bool) {
			result, err := libtime.ParseUnixTime(ctx, input)
			if expectedError {
				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Format(time.RFC3339)).To(Equal(expectedDateString))
			}
		},
		Entry("invalid", "banana", "", true),
		Entry("dateTime", "2023-06-19T07:56:34Z", "2023-06-19T07:56:34Z", false),
		Entry("unixTime", 1687161394, "2023-06-19T07:56:34Z", false),
		Entry("unixTimeStr", "1687161394", "2023-06-19T07:56:34Z", false),
	)
	Context("TimePtr", func() {
		var dateTime *libtime.UnixTime
		var timePtr *time.Time
		BeforeEach(func() {
			dateTime = libtime.UnixTime(time.Unix(1000, 0)).Ptr()
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
		var dateTime libtime.UnixTime
		var result libtime.UnixTime
		var days int
		var months int
		var years int
		BeforeEach(func() {
			years = 0
			months = 0
			days = 0
			dateTime = ParseUnixTime("2024-12-24T20:15:59Z")
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
				var unixTime libtime.UnixTime
				Expect(unixTime.IsZero()).To(BeTrue())
			})
		})
		Context("non-zero time", func() {
			It("returns false", func() {
				unixTime := libtime.UnixTimeFromSeconds(1687161394)
				Expect(unixTime.IsZero()).To(BeFalse())
			})
		})
	})
	Context("NewUnixTime", func() {
		var result libtime.UnixTime
		BeforeEach(func() {
			result = libtime.NewUnixTime(2023, time.June, 19, 7, 56, 34, 0, time.UTC)
		})
		It("creates correct UnixTime", func() {
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
			UnixTime    libtime.UnixTime  `yaml:"unixTime"`
			UnixTimePtr *libtime.UnixTime `yaml:"unixTimePtr"`
		}
		var original TestStruct
		var unmarshaled TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			original = TestStruct{
				UnixTime:    libtime.UnixTime(time.Unix(1687161394, 0)),
				UnixTimePtr: libtime.UnixTime(time.Unix(1687161394, 0)).Ptr(),
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
		It("marshals to RFC3339Nano format (not integer)", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`unixTime: "?2023-06-19T07:56:34Z"?`))
			Expect(yamlString).To(MatchRegexp(`unixTimePtr: "?2023-06-19T07:56:34Z"?`))
		})
		It("round-trips UnixTime field correctly", func() {
			Expect(unmarshaled.UnixTime.String()).To(Equal("2023-06-19T07:56:34Z"))
			Expect(unmarshaled.UnixTime.String()).To(Equal(original.UnixTime.String()))
		})
		It("round-trips UnixTimePtr field correctly", func() {
			Expect(unmarshaled.UnixTimePtr).NotTo(BeNil())
			Expect(unmarshaled.UnixTimePtr.String()).To(Equal("2023-06-19T07:56:34Z"))
			Expect(unmarshaled.UnixTimePtr.String()).To(Equal(original.UnixTimePtr.String()))
		})
	})
	Context("YAML omitempty - Phase 2", func() {
		type TestStruct struct {
			UnixTime        libtime.UnixTime  `yaml:"unixTime,omitempty"`
			UnixTimePtr     *libtime.UnixTime `yaml:"unixTimePtr,omitempty"`
			UnixTimeNonZero libtime.UnixTime  `yaml:"unixTimeNonZero,omitempty"`
		}
		var testStruct TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			testStruct = TestStruct{
				UnixTime:        libtime.UnixTime{},                         // zero
				UnixTimePtr:     nil,                                        // nil
				UnixTimeNonZero: libtime.UnixTime(time.Unix(1687161394, 0)), // non-zero
			}
		})
		JustBeforeEach(func() {
			yamlBytes, err = yaml.Marshal(testStruct)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("omits zero UnixTime field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("unixTime:"))
		})
		It("omits nil *UnixTime field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("unixTimePtr:"))
		})
		It("includes non-zero UnixTime field", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`unixTimeNonZero: "?2023-06-19T07:56:34Z"?`))
		})
	})
	Context("struct marshal regression - Phase 1", func() {
		Context("A. JSON Marshal — non-zero values set for Field and FieldPtr only", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.UnixTime  `json:"field"`
					FieldPtr     *libtime.UnixTime `json:"fieldPtr"`
					FieldOmit    libtime.UnixTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.UnixTime `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    libtime.UnixTime(time.Unix(1687161394, 0)),
					FieldPtr: libtime.UnixTime(time.Unix(1687161394, 0)).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns exact JSON output", func() {
				Expect(
					string(jsonBytes),
				).To(Equal(`{"field":1687161394,"fieldPtr":1687161394,"fieldOmit":-62135596800}`))
			})
		})
		Context("B. JSON Unmarshal — round-trip", func() {
			var original struct {
				Field        libtime.UnixTime  `json:"field"`
				FieldPtr     *libtime.UnixTime `json:"fieldPtr"`
				FieldOmit    libtime.UnixTime  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.UnixTime `json:"fieldPtrOmit,omitempty"`
			}
			var unmarshaled struct {
				Field        libtime.UnixTime  `json:"field"`
				FieldPtr     *libtime.UnixTime `json:"fieldPtr"`
				FieldOmit    libtime.UnixTime  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.UnixTime `json:"fieldPtrOmit,omitempty"`
			}
			var jsonBytes []byte
			BeforeEach(func() {
				original = struct {
					Field        libtime.UnixTime  `json:"field"`
					FieldPtr     *libtime.UnixTime `json:"fieldPtr"`
					FieldOmit    libtime.UnixTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.UnixTime `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    libtime.UnixTime(time.Unix(1687161394, 0)),
					FieldPtr: libtime.UnixTime(time.Unix(1687161394, 0)).Ptr(),
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
				Expect(unmarshaled.Field.Unix()).To(Equal(original.Field.Unix()))
			})
			It("round-trips FieldPtr correctly", func() {
				Expect(unmarshaled.FieldPtr).NotTo(BeNil())
				Expect(unmarshaled.FieldPtr.Unix()).To(Equal(original.FieldPtr.Unix()))
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
					Field        libtime.UnixTime  `json:"field"`
					FieldPtr     *libtime.UnixTime `json:"fieldPtr"`
					FieldOmit    libtime.UnixTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.UnixTime `json:"fieldPtrOmit,omitempty"`
				}{
					Field:        libtime.UnixTime(time.Unix(1687161394, 0)),
					FieldPtr:     libtime.UnixTime(time.Unix(1687161394, 0)).Ptr(),
					FieldOmit:    libtime.UnixTime(time.Unix(1687161394, 0)),
					FieldPtrOmit: libtime.UnixTime(time.Unix(1687161394, 0)).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("all fields appear in output", func() {
				jsonStr := string(jsonBytes)
				Expect(jsonStr).To(ContainSubstring(`"field":1687161394`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtr":1687161394`))
				Expect(jsonStr).To(ContainSubstring(`"fieldOmit":1687161394`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtrOmit":1687161394`))
			})
		})
		Context("D. JSON Marshal — all fields zero/nil", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.UnixTime  `json:"field"`
					FieldPtr     *libtime.UnixTime `json:"fieldPtr"`
					FieldOmit    libtime.UnixTime  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.UnixTime `json:"fieldPtrOmit,omitempty"`
				}{}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It(
				"non-omitempty fields produce -62135596800, omitempty fields produce -62135596800",
				func() {
					Expect(
						string(jsonBytes),
					).To(Equal(`{"field":-62135596800,"fieldPtr":null,"fieldOmit":-62135596800}`))
				},
			)
		})
	})
})
