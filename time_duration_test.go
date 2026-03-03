// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time_test

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	libtime "github.com/bborbe/time"
)

var _ = DescribeTable("ParseDuration",
	func(input string, expectedDuration libtime.Duration, expectedError bool) {
		duration, err := libtime.ParseDuration(context.Background(), input)
		if expectedError {
			Expect(err).NotTo(BeNil())
			Expect(duration).To(BeNil())
		} else {
			Expect(err).To(BeNil())
			Expect(duration).NotTo(BeNil())
			Expect(*duration).To(Equal(expectedDuration))
		}
	},
	Entry("without unit", "1337", libtime.Duration(1337), false),
	Entry("ns", "1ns", libtime.Nanosecond, false),
	Entry("us", "1us", libtime.Microsecond, false),
	Entry("ms", "1ms", libtime.Millisecond, false),
	Entry("second", "1s", libtime.Second, false),
	Entry("minute", "1m", libtime.Minute, false),
	Entry("hour", "1h", libtime.Hour, false),
	Entry("day", "1d", 24*libtime.Hour, false),
	Entry("week", "1w", 7*24*libtime.Hour, false),
	// Uppercase variants
	Entry("second uppercase", "1S", libtime.Second, false),
	Entry("minute uppercase", "1M", libtime.Minute, false),
	Entry("hour uppercase", "1H", libtime.Hour, false),
	Entry("day uppercase", "1D", 24*libtime.Hour, false),
	Entry("week uppercase", "1W", 7*24*libtime.Hour, false),
	// Mixed case combinations
	Entry("combined uppercase", "1H30M", 90*libtime.Minute, false),
	Entry("combined mixed case", "1h30M", 90*libtime.Minute, false),
	Entry("combined mixed case 2", "1H30m", 90*libtime.Minute, false),
	Entry("negative uppercase", "-1H30M", -90*libtime.Minute, false),
	Entry("dot uppercase", "1.5H", 90*libtime.Minute, false),
	// Positive prefix
	Entry("positive prefix", "+1h", libtime.Hour, false),
	Entry("positive prefix combined", "+1h30m", 90*libtime.Minute, false),
	// Original lowercase tests
	Entry("combined", "1h30m", 90*libtime.Minute, false),
	Entry("negative", "-1h30m", -90*libtime.Minute, false),
	Entry("dot", "1.5h", 90*libtime.Minute, false),
	Entry("hello", "hello", libtime.Duration(0), true),
	Entry("hello1d", "hello1d", libtime.Duration(0), true),
)

var _ = Describe("Duration", func() {
	var err error
	var _ = DescribeTable(
		"String",
		func(inputDuration libtime.Duration, expectedOutput string) {
			Expect(inputDuration.String()).To(Equal(expectedOutput))
		},
		Entry("30s", 30*libtime.Second, "30s"),
		Entry("59m30s", 59*libtime.Minute+30*libtime.Second, "59m30s"),
		Entry("23h59m30s", 23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "23h59m30s"),
		Entry(
			"5d23h59m30s",
			5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second,
			"5d23h59m30s",
		),
		Entry(
			"10w5d23h59m30s",
			10*libtime.Week+5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second,
			"10w5d23h59m30s",
		),
	)

	var _ = DescribeTable(
		"MarshalJSON",
		func(inputDuration libtime.Duration, expectedOutput string, expectError bool) {
			bytes, err := inputDuration.MarshalJSON()
			if expectError {
				Expect(err).NotTo(BeNil())
				Expect(bytes).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(string(bytes)).To(Equal(expectedOutput))
			}
		},
		Entry("0", libtime.Duration(0), `"0s"`, false),
		Entry("30s", 30*libtime.Second, `"30s"`, false),
		Entry("59m30s", 59*libtime.Minute+30*libtime.Second, `"59m30s"`, false),
		Entry(
			"23h59m30s",
			23*libtime.Hour+59*libtime.Minute+30*libtime.Second,
			`"23h59m30s"`,
			false,
		),
		Entry(
			"143h59m30s",
			5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second,
			`"143h59m30s"`,
			false,
		),
		Entry(
			"1823h59m30s",
			10*libtime.Week+5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second,
			`"1823h59m30s"`,
			false,
		),
	)

	var _ = DescribeTable(
		"String",
		func(inputDuration libtime.Duration, expectedOutput string) {
			Expect(inputDuration.String()).To(Equal(expectedOutput))
		},
		Entry("1w", libtime.Week, "1w"),
		Entry("1d", libtime.Day, "1d"),
		Entry("1h", libtime.Hour, "1h"),
		Entry("1m", libtime.Minute, "1m"),
		Entry("1s", libtime.Second, "1s"),
		Entry("1ms", libtime.Millisecond, "1ms"),
		Entry("1µs", libtime.Microsecond, "1µs"),
		Entry("1ns", libtime.Nanosecond, "1ns"),
		Entry("0", libtime.Duration(0), "0s"),
		Entry("1w1ns", libtime.Week+libtime.Nanosecond, "1w1ns"),
		Entry("59m30s", 59*libtime.Minute+30*libtime.Second, "59m30s"),
		Entry("23h59m30s", 23*libtime.Hour+59*libtime.Minute+30*libtime.Second, "23h59m30s"),
		Entry(
			"5d23h59m30s",
			5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second,
			"5d23h59m30s",
		),
		Entry(
			"10w5d23h59m30s",
			10*libtime.Week+5*libtime.Day+23*libtime.Hour+59*libtime.Minute+30*libtime.Second,
			"10w5d23h59m30s",
		),
	)

	Context("UnmarshalJSON", func() {
		var err error
		var duration libtime.Duration
		var value string
		BeforeEach(func() {
			duration = 0
		})
		JustBeforeEach(func() {
			err = duration.UnmarshalJSON([]byte(value))
		})
		Context("with string value", func() {
			BeforeEach(func() {
				value = `"1337"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(libtime.Duration(1337)))
			})
		})
		Context("with number value", func() {
			BeforeEach(func() {
				value = `1337`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(libtime.Duration(1337)))
			})
		})
		Context("with duration value", func() {
			BeforeEach(func() {
				value = `"1h"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(libtime.Hour))
			})
		})
		Context("with uppercase duration value", func() {
			BeforeEach(func() {
				value = `"1H"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(libtime.Hour))
			})
		})
		Context("with mixed case duration value", func() {
			BeforeEach(func() {
				value = `"1H30m"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns correct content", func() {
				Expect(duration).To(Equal(90 * libtime.Minute))
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
				Expect(duration).To(Equal(libtime.Duration(0)))
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
				Expect(duration).To(Equal(libtime.Duration(0)))
			})
		})
	})
	Context("YAML round-trip - Phase 2", func() {
		type TestStruct struct {
			Duration    libtime.Duration  `yaml:"duration"`
			DurationPtr *libtime.Duration `yaml:"durationPtr"`
		}
		var original TestStruct
		var unmarshaled TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			original = TestStruct{
				Duration:    90 * libtime.Minute,
				DurationPtr: (90 * libtime.Minute).Ptr(),
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
		It("marshals to Go duration format", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`duration: "?1h30m0s"?`))
			Expect(yamlString).To(MatchRegexp(`durationPtr: "?1h30m0s"?`))
		})
		It("round-trips Duration field correctly", func() {
			Expect(unmarshaled.Duration).To(Equal(original.Duration))
		})
		It("round-trips DurationPtr field correctly", func() {
			Expect(unmarshaled.DurationPtr).NotTo(BeNil())
			Expect(*unmarshaled.DurationPtr).To(Equal(*original.DurationPtr))
		})
	})
	Context("YAML omitempty - Phase 2", func() {
		type TestStruct struct {
			Duration        libtime.Duration  `yaml:"duration,omitempty"`
			DurationPtr     *libtime.Duration `yaml:"durationPtr,omitempty"`
			DurationNonZero libtime.Duration  `yaml:"durationNonZero,omitempty"`
		}
		var testStruct TestStruct
		var yamlBytes []byte
		BeforeEach(func() {
			testStruct = TestStruct{
				Duration:        libtime.Duration(0), // zero
				DurationPtr:     nil,                 // nil
				DurationNonZero: 90 * libtime.Minute, // non-zero
			}
		})
		JustBeforeEach(func() {
			yamlBytes, err = yaml.Marshal(testStruct)
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("omits zero Duration field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("duration:"))
		})
		It("omits nil *Duration field with omitempty", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).NotTo(ContainSubstring("durationPtr:"))
		})
		It("includes non-zero Duration field", func() {
			yamlString := string(yamlBytes)
			Expect(yamlString).To(MatchRegexp(`durationNonZero: "?1h30m0s"?`))
		})
	})
	Context("struct marshal regression - Phase 1", func() {
		Context("A. JSON Marshal — non-zero values set for Field and FieldPtr only", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.Duration  `json:"field"`
					FieldPtr     *libtime.Duration `json:"fieldPtr"`
					FieldOmit    libtime.Duration  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Duration `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    90 * libtime.Minute,
					FieldPtr: (90 * libtime.Minute).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns exact JSON output", func() {
				Expect(string(jsonBytes)).To(Equal(`{"field":"1h30m0s","fieldPtr":"1h30m0s"}`))
			})
		})
		Context("B. JSON Unmarshal — round-trip", func() {
			var original struct {
				Field        libtime.Duration  `json:"field"`
				FieldPtr     *libtime.Duration `json:"fieldPtr"`
				FieldOmit    libtime.Duration  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.Duration `json:"fieldPtrOmit,omitempty"`
			}
			var unmarshaled struct {
				Field        libtime.Duration  `json:"field"`
				FieldPtr     *libtime.Duration `json:"fieldPtr"`
				FieldOmit    libtime.Duration  `json:"fieldOmit,omitempty"`
				FieldPtrOmit *libtime.Duration `json:"fieldPtrOmit,omitempty"`
			}
			var jsonBytes []byte
			BeforeEach(func() {
				original = struct {
					Field        libtime.Duration  `json:"field"`
					FieldPtr     *libtime.Duration `json:"fieldPtr"`
					FieldOmit    libtime.Duration  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Duration `json:"fieldPtrOmit,omitempty"`
				}{
					Field:    90 * libtime.Minute,
					FieldPtr: (90 * libtime.Minute).Ptr(),
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
				Expect(unmarshaled.Field).To(Equal(original.Field))
			})
			It("round-trips FieldPtr correctly", func() {
				Expect(unmarshaled.FieldPtr).NotTo(BeNil())
				Expect(*unmarshaled.FieldPtr).To(Equal(*original.FieldPtr))
			})
			It("zero fields remain zero", func() {
				Expect(unmarshaled.FieldOmit).To(Equal(libtime.Duration(0)))
				Expect(unmarshaled.FieldPtrOmit).To(BeNil())
			})
		})
		Context("C. JSON Marshal — all fields set (non-zero)", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.Duration  `json:"field"`
					FieldPtr     *libtime.Duration `json:"fieldPtr"`
					FieldOmit    libtime.Duration  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Duration `json:"fieldPtrOmit,omitempty"`
				}{
					Field:        90 * libtime.Minute,
					FieldPtr:     (90 * libtime.Minute).Ptr(),
					FieldOmit:    90 * libtime.Minute,
					FieldPtrOmit: (90 * libtime.Minute).Ptr(),
				}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("all fields appear in output", func() {
				jsonStr := string(jsonBytes)
				Expect(jsonStr).To(ContainSubstring(`"field":"1h30m0s"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtr":"1h30m0s"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldOmit":"1h30m0s"`))
				Expect(jsonStr).To(ContainSubstring(`"fieldPtrOmit":"1h30m0s"`))
			})
		})
		Context("D. JSON Marshal — all fields zero/nil", func() {
			var jsonBytes []byte
			JustBeforeEach(func() {
				testStruct := struct {
					Field        libtime.Duration  `json:"field"`
					FieldPtr     *libtime.Duration `json:"fieldPtr"`
					FieldOmit    libtime.Duration  `json:"fieldOmit,omitempty"`
					FieldPtrOmit *libtime.Duration `json:"fieldPtrOmit,omitempty"`
				}{}
				jsonBytes, err = json.Marshal(testStruct)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("non-omitempty fields produce 0s, omitempty fields omitted", func() {
				Expect(string(jsonBytes)).To(Equal(`{"field":"0s","fieldPtr":null}`))
			})
		})
	})
})
