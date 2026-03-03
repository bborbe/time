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

var _ = Describe("TimeOfDay", func() {
	var err error
	var timeOfDay libtime.TimeOfDay
	var now time.Time
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
		now = ParseTime("2023-05-02T12:45:59.123456Z")
	})
	JustBeforeEach(func() {
		libtime.Now = func() time.Time {
			return now
		}
	})
	Context("ParseTimeOfDay", func() {
		var input string
		var timeOfDay *libtime.TimeOfDay
		var winterTime libtime.DateTime
		var summerTime libtime.DateTime
		JustBeforeEach(func() {
			timeOfDay, err = libtime.ParseTimeOfDay(ctx, input)
			Expect(err).To(BeNil())
			winterTime = timeOfDay.DateTime(2024, time.January, 1)
			summerTime = timeOfDay.DateTime(2024, time.July, 1)
		})
		Context("NOW", func() {
			BeforeEach(func() {
				input = "NOW"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T12:45:59.123456Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T12:45:59.123456Z"))
			})
		})
		Context("time with Z", func() {
			BeforeEach(func() {
				input = "13:37:59.123456Z"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T13:37:59.123456Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T13:37:59.123456Z"))
			})
		})
		Context("time with UTC", func() {
			BeforeEach(func() {
				input = "14:37:59 UTC"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T14:37:59Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T14:37:59Z"))
			})
		})
		Context("time with Europe/Berlin", func() {
			BeforeEach(func() {
				input = "15:37:59 Europe/Berlin"
			})
			It("returns correct winterTime", func() {
				Expect(winterTime).NotTo(BeNil())
				Expect(winterTime.Format(time.RFC3339Nano)).To(Equal("2024-01-01T15:37:59+01:00"))
				Expect(
					winterTime.Time().UTC().Format(time.RFC3339Nano),
				).To(Equal("2024-01-01T14:37:59Z"))
			})
			It("returns correct summerTime", func() {
				Expect(summerTime).NotTo(BeNil())
				Expect(summerTime.Format(time.RFC3339Nano)).To(Equal("2024-07-01T15:37:59+02:00"))
				Expect(
					summerTime.Time().UTC().Format(time.RFC3339Nano),
				).To(Equal("2024-07-01T13:37:59Z"))
			})
		})
	})
	Context("String", func() {
		var result string
		JustBeforeEach(func() {
			result = timeOfDay.String()
		})
		Context("with nano", func() {
			BeforeEach(func() {
				timeOfDay = libtime.TimeOfDayFromTime(ParseTime("2023-05-02T13:45:59.123456Z"))
			})
			It("returns correct string", func() {
				Expect(result).To(Equal("13:45:59.123456Z"))
			})
		})
		Context("without nano", func() {
			BeforeEach(func() {
				timeOfDay = libtime.TimeOfDayFromTime(ParseTime("2023-05-02T13:45:59Z"))
			})
			It("returns correct string", func() {
				Expect(result).To(Equal("13:45:59Z"))
			})
		})
	})
	Context("MarshalJSON", func() {
		var bytes []byte
		BeforeEach(func() {
			timeOfDay = libtime.TimeOfDayFromTime(ParseTime("2023-05-02T13:45:59.123456Z"))
		})
		JustBeforeEach(func() {
			bytes, err = timeOfDay.MarshalJSON()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns correct content", func() {
			Expect(string(bytes)).To(Equal(`"13:45:59.123456Z"`))
		})
	})
	DescribeTable(
		"Date",
		func(input libtime.TimeOfDay, year int, month int, day int, expectedTime string, expectError bool) {
			dateTime, err := input.Date(year, time.Month(month), day)
			if expectError {
				Expect(err).NotTo(BeNil())
				Expect(timeOfDay).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(dateTime.UTC().Format(time.RFC3339)).To(Equal(expectedTime))
			}
		},
		Entry("13:37", ParseTimeOfDay("13:37"), 2024, 12, 24, "2024-12-24T13:37:00Z", false),
		Entry("13:37:42", ParseTimeOfDay("13:37:42"), 2024, 12, 24, "2024-12-24T13:37:42Z", false),
		Entry(
			"13:37:42Z",
			ParseTimeOfDay("13:37:42Z"),
			2024,
			12,
			24,
			"2024-12-24T13:37:42Z",
			false,
		),
		Entry(
			"13:37:42 Europe/Berlin",
			ParseTimeOfDay("13:37:42 Europe/Berlin"),
			2024,
			12,
			24,
			"2024-12-24T12:37:42Z",
			false,
		),
	)
	DescribeTable(
		"UnmarshalJSON",
		func(input string, expected string, expectError bool) {
			timeOfDay = libtime.TimeOfDay{}
			err = timeOfDay.UnmarshalJSON([]byte(input))
			if expectError {
				Expect(err).NotTo(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(timeOfDay.String()).To(Equal(expected))
			}
		},
		Entry("error", `"banana"`, ``, true),
		Entry("with tz", `"13:45:59Z"`, `13:45:59Z`, false),
		Entry("with tz and ns", `"13:45:59.123456Z"`, `13:45:59.123456Z`, false),
		Entry("hour:min tz", `"13:45Z"`, `13:45:00Z`, false),
		Entry("hour:min", `"13:45"`, `13:45:00Z`, false),
		Entry("without tz", `"13:45:59"`, `13:45:59Z`, false),
		Entry("without tz and ns", `"13:45:59.123456"`, `13:45:59.123456Z`, false),
		Entry("datetime with tz", `"2023-10-02T13:45:59Z"`, `13:45:59Z`, false),
		Entry(
			"datetime with tz and ns",
			`"2023-10-02T13:45:59.123456Z"`,
			`13:45:59.123456Z`,
			false,
		),
	)
	Context("YAML round-trip - Phase 2", func() {
		type TestStruct struct {
			TimeOfDay    libtime.TimeOfDay  `yaml:"timeOfDay"`
			TimeOfDayPtr *libtime.TimeOfDay `yaml:"timeOfDayPtr"`
		}
		var original TestStruct
		var unmarshaled TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			original = TestStruct{
				TimeOfDay:    libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC},
				TimeOfDayPtr: (&libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC}),
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
		It("marshals to time format", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`timeOfDay: "?13:37:00Z"?`))
			Expect(yamlString).To(MatchRegexp(`timeOfDayPtr: "?13:37:00Z"?`))
		})
		It("round-trips TimeOfDay field correctly", func() {
			Expect(unmarshaled.TimeOfDay.Equal(original.TimeOfDay)).To(BeTrue())
		})
		It("round-trips TimeOfDayPtr field correctly", func() {
			Expect(unmarshaled.TimeOfDayPtr).NotTo(BeNil())
			Expect(unmarshaled.TimeOfDayPtr.Equal(*original.TimeOfDayPtr)).To(BeTrue())
		})
	})
	Context("YAML omitempty - Phase 2", func() {
		type TestStruct struct {
			TimeOfDay        libtime.TimeOfDay  `yaml:"timeOfDay,omitempty"`
			TimeOfDayPtr     *libtime.TimeOfDay `yaml:"timeOfDayPtr,omitempty"`
			TimeOfDayNonZero libtime.TimeOfDay  `yaml:"timeOfDayNonZero,omitempty"`
		}
		var testStruct TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			testStruct = TestStruct{
				TimeOfDay: libtime.TimeOfDay{
					Location: time.UTC,
				}, // zero
				TimeOfDayPtr: nil, // nil
				TimeOfDayNonZero: libtime.TimeOfDay{
					Hour:     13,
					Minute:   37,
					Location: time.UTC,
				}, // non-zero
			}
		})
		JustBeforeEach(func() {
			yamlBytes, err = yaml.Marshal(testStruct)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("includes zero TimeOfDay field with omitempty (YAML behavior)", func() {
			yamlString := string(yamlBytes)
			// YAML may include zero TimeOfDay as "00:00:00Z"
			Expect(yamlString).To(MatchRegexp(`timeOfDay: "?00:00:00Z"?`))
		})
		It("omits nil *TimeOfDay field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("timeOfDayPtr:"))
		})
		It("includes non-zero TimeOfDay field", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`timeOfDayNonZero: "?13:37:00Z"?`))
		})
	})
	Context("struct marshal regression - Phase 1", func() {
		Context("A. JSON Marshal — non-zero values set for Field and FieldPtr only", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.TimeOfDay  `json:"field"`
					FieldPtr     *libtime.TimeOfDay `json:"fieldPtr"`
					FieldOmit    libtime.TimeOfDay  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.TimeOfDay `json:"fieldPtrOmit,omitempty"`
				}{
					Field:     libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC},
					FieldPtr:  (&libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC}),
					FieldOmit: libtime.TimeOfDay{Location: time.UTC},
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns exact JSON output", func() {
				Expect(
					string(jsonBytes),
				).To(Equal(`{"field":"13:37:00Z","fieldPtr":"13:37:00Z","fieldOmit":"00:00:00Z"}`))
			})
		})
		Context("B. JSON Unmarshal — round-trip", func() {
			var original struct {
				Field        libtime.TimeOfDay  `json:"field"`
				FieldPtr     *libtime.TimeOfDay `json:"fieldPtr"`
				FieldOmit    libtime.TimeOfDay  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.TimeOfDay `json:"fieldPtrOmit,omitempty"`
			}
			var unmarshaled struct {
				Field        libtime.TimeOfDay  `json:"field"`
				FieldPtr     *libtime.TimeOfDay `json:"fieldPtr"`
				FieldOmit    libtime.TimeOfDay  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.TimeOfDay `json:"fieldPtrOmit,omitempty"`
			}
			var jsonBytes []byte
			BeforeEach(func() {
				original = struct {
					Field        libtime.TimeOfDay  `json:"field"`
					FieldPtr     *libtime.TimeOfDay `json:"fieldPtr"`
					FieldOmit    libtime.TimeOfDay  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.TimeOfDay `json:"fieldPtrOmit,omitempty"`
				}{
					Field:     libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC},
					FieldPtr:  (&libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC}),
					FieldOmit: libtime.TimeOfDay{Location: time.UTC},
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
				Expect(unmarshaled.Field.Equal(original.Field)).To(BeTrue())
			})
			It("round-trips FieldPtr correctly", func() {
				Expect(unmarshaled.FieldPtr).NotTo(BeNil())
				Expect(unmarshaled.FieldPtr.Equal(*original.FieldPtr)).To(BeTrue())
			})
			It("zero fields remain zero", func() {
				Expect(unmarshaled.FieldOmit.Hour).To(Equal(0))
				Expect(unmarshaled.FieldOmit.Minute).To(Equal(0))
				Expect(unmarshaled.FieldPtrOmit).To(BeNil())
			})
		})
		Context("C. JSON Marshal — all fields set (non-zero)", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.TimeOfDay  `json:"field"`
					FieldPtr     *libtime.TimeOfDay `json:"fieldPtr"`
					FieldOmit    libtime.TimeOfDay  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.TimeOfDay `json:"fieldPtrOmit,omitempty"`
				}{
					Field:        libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC},
					FieldPtr:     (&libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC}),
					FieldOmit:    libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC},
					FieldPtrOmit: (&libtime.TimeOfDay{Hour: 13, Minute: 37, Location: time.UTC}),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("all fields appear in output", func() {
				jsonStr := string(jsonBytes)
				Expect(jsonStr).To(ContainSubstring(`"field":"13:37:00Z"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtr":"13:37:00Z"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldOmit":"13:37:00Z"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtrOmit":"13:37:00Z"`))
			})
		})
		Context("D. JSON Marshal — all fields zero/nil", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.TimeOfDay  `json:"field"`
					FieldPtr     *libtime.TimeOfDay `json:"fieldPtr"`
					FieldOmit    libtime.TimeOfDay  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.TimeOfDay `json:"fieldPtrOmit,omitempty"`
				}{
					Field:     libtime.TimeOfDay{Location: time.UTC},
					FieldOmit: libtime.TimeOfDay{Location: time.UTC},
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It(
				"non-omitempty fields produce 00:00:00Z, omitempty fields produce 00:00:00Z",
				func() {
					Expect(
						string(jsonBytes),
					).To(Equal(`{"field":"00:00:00Z","fieldPtr":null,"fieldOmit":"00:00:00Z"}`))
				},
			)
		})
	})
})
