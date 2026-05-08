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

	libtime "github.com/bborbe/time"
)

var _ = Describe("DateOrDateTime", func() {
	var err error
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("MarshalBinary & DateOrDateTimeFromBinary", func() {
		var value libtime.DateOrDateTime
		var binary []byte
		BeforeEach(func() {
			value = libtime.DateOrDateTime(time.Unix(1687161394, 0).UTC())
		})
		JustBeforeEach(func() {
			binary, err = value.MarshalBinary()
		})
		It("returns no error", func() {
			Expect(err).To(BeNil())
		})
		It("returns binary", func() {
			Expect(binary).NotTo(BeNil())
		})
		It("round-trips via DateOrDateTimeFromBinary", func() {
			result, err := libtime.DateOrDateTimeFromBinary(ctx, binary)
			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Unix()).To(Equal(int64(1687161394)))
		})
	})

	Context("MarshalJSON", func() {
		var value libtime.DateOrDateTime
		var jsonBytes []byte
		JustBeforeEach(func() {
			jsonBytes, err = value.MarshalJSON()
		})
		Context("midnight UTC", func() {
			BeforeEach(func() {
				value = libtime.DateOrDateTime(
					time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
				)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("produces date-only format", func() {
				Expect(string(jsonBytes)).To(Equal(`"2026-01-15"`))
			})
		})
		Context("non-midnight UTC", func() {
			BeforeEach(func() {
				value = libtime.DateOrDateTime(
					time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC),
				)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("produces RFC3339Nano format", func() {
				Expect(string(jsonBytes)).To(Equal(`"2026-01-15T14:30:00Z"`))
			})
		})
		Context("zero value", func() {
			BeforeEach(func() {
				value = libtime.DateOrDateTime{}
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("produces null", func() {
				Expect(string(jsonBytes)).To(Equal(`null`))
			})
		})
	})

	Context("UnmarshalJSON", func() {
		var value libtime.DateOrDateTime
		var input string
		JustBeforeEach(func() {
			err = value.UnmarshalJSON([]byte(input))
		})
		Context("date-only input", func() {
			BeforeEach(func() {
				input = `"2026-01-15"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("re-marshals as date-only", func() {
				b, e := value.MarshalJSON()
				Expect(e).To(BeNil())
				Expect(string(b)).To(Equal(`"2026-01-15"`))
			})
		})
		Context("RFC3339 input", func() {
			BeforeEach(func() {
				input = `"2026-01-15T14:30:00Z"`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("re-marshals as RFC3339Nano", func() {
				b, e := value.MarshalJSON()
				Expect(e).To(BeNil())
				Expect(string(b)).To(Equal(`"2026-01-15T14:30:00Z"`))
			})
		})
		Context("empty string", func() {
			BeforeEach(func() {
				input = `""`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("yields zero value", func() {
				Expect(value.IsZero()).To(BeTrue())
			})
		})
		Context("null", func() {
			BeforeEach(func() {
				input = `null`
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("yields zero value", func() {
				Expect(value.IsZero()).To(BeTrue())
			})
		})
		Context("invalid input", func() {
			BeforeEach(func() {
				input = `"not-a-date"`
			})
			It("returns non-nil error", func() {
				Expect(err).NotTo(BeNil())
			})
		})
	})

	Context("MarshalText", func() {
		var value libtime.DateOrDateTime
		var textBytes []byte
		JustBeforeEach(func() {
			textBytes, err = value.MarshalText()
		})
		Context("midnight UTC", func() {
			BeforeEach(func() {
				value = libtime.DateOrDateTime(
					time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
				)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("produces date-only bytes", func() {
				Expect(string(textBytes)).To(Equal("2026-01-15"))
			})
		})
		Context("non-midnight", func() {
			BeforeEach(func() {
				value = libtime.DateOrDateTime(
					time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC),
				)
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("produces RFC3339Nano bytes", func() {
				Expect(string(textBytes)).To(Equal("2026-01-15T14:30:00Z"))
			})
		})
		Context("midnight in non-UTC zone (not midnight UTC)", func() {
			BeforeEach(func() {
				loc := time.FixedZone("UTC+2", 2*60*60)
				value = libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, loc))
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("produces RFC3339Nano bytes (not date-only)", func() {
				result := string(textBytes)
				Expect(result).NotTo(Equal("2026-01-15"))
				Expect(result).To(ContainSubstring("T"))
			})
		})
		Context("zero value", func() {
			BeforeEach(func() {
				value = libtime.DateOrDateTime{}
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns nil bytes", func() {
				Expect(textBytes).To(BeNil())
			})
		})
	})

	Context("UnmarshalText", func() {
		var value libtime.DateOrDateTime
		var textInput []byte
		JustBeforeEach(func() {
			err = value.UnmarshalText(textInput)
		})
		Context("date-only input", func() {
			BeforeEach(func() {
				textInput = []byte("2026-01-15")
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("round-trips as date-only", func() {
				Expect(value.String()).To(Equal("2026-01-15"))
			})
		})
		Context("RFC3339 input", func() {
			BeforeEach(func() {
				textInput = []byte("2026-01-15T14:30:00Z")
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("round-trips as RFC3339Nano", func() {
				Expect(value.String()).To(Equal("2026-01-15T14:30:00Z"))
			})
		})
		Context("empty bytes", func() {
			BeforeEach(func() {
				textInput = []byte("")
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("yields zero value", func() {
				Expect(value.IsZero()).To(BeTrue())
			})
		})
	})

	Context("IsDateOnly", func() {
		It("returns true for midnight UTC", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			Expect(v.IsDateOnly()).To(BeTrue())
		})
		It("returns false for non-midnight", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC))
			Expect(v.IsDateOnly()).To(BeFalse())
		})
		It("returns false for zero value", func() {
			var v libtime.DateOrDateTime
			Expect(v.IsDateOnly()).To(BeFalse())
		})
	})

	Context("IsZero", func() {
		It("returns true for zero value", func() {
			var v libtime.DateOrDateTime
			Expect(v.IsZero()).To(BeTrue())
		})
		It("returns false for non-zero value", func() {
			v := libtime.DateOrDateTime(time.Unix(1687161394, 0))
			Expect(v.IsZero()).To(BeFalse())
		})
	})

	Context("AsDate", func() {
		It("converts midnight UTC to the correct Date", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			d := v.AsDate()
			Expect(d.String()).To(Equal("2026-01-15"))
		})
	})

	Context("AsDateTime", func() {
		It("converts to DateTime with matching time", func() {
			t := time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC)
			v := libtime.DateOrDateTime(t)
			dt := v.AsDateTime()
			Expect(dt.Time().Equal(t)).To(BeTrue())
		})
	})

	Context("String", func() {
		It("returns YYYY-MM-DD for midnight UTC", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			Expect(v.String()).To(Equal("2026-01-15"))
		})
		It("returns RFC3339Nano for non-midnight", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC))
			Expect(v.String()).To(Equal("2026-01-15T14:30:00Z"))
		})
		It("returns empty string for zero value", func() {
			var v libtime.DateOrDateTime
			Expect(v.String()).To(Equal(""))
		})
	})

	Context("Validate", func() {
		It("returns error for zero value", func() {
			var v libtime.DateOrDateTime
			Expect(v.Validate(ctx)).NotTo(BeNil())
		})
		It("returns nil for non-zero value", func() {
			v := libtime.DateOrDateTime(time.Unix(1687161394, 0))
			Expect(v.Validate(ctx)).To(BeNil())
		})
	})

	Context("UTC", func() {
		It("converts non-UTC zone to UTC", func() {
			loc := time.FixedZone("UTC+5", 5*60*60)
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 5, 0, 0, 0, loc))
			utc := v.UTC()
			Expect(utc.Time().Location()).To(Equal(time.UTC))
			Expect(utc.Time().Hour()).To(Equal(0))
		})
	})

	Context("Clone", func() {
		It("returns equal value for non-zero value", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			Expect(v.Clone()).To(Equal(v))
		})
		It("returns zero value for zero value", func() {
			var v libtime.DateOrDateTime
			Expect(v.Clone()).To(Equal(libtime.DateOrDateTime{}))
		})
	})

	Context("ClonePtr", func() {
		It("returns nil for nil pointer", func() {
			var v *libtime.DateOrDateTime
			Expect(v.ClonePtr()).To(BeNil())
		})
		It("returns non-nil pointer equal to original", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			result := v.ClonePtr()
			Expect(result).NotTo(BeNil())
			Expect(*result).To(Equal(v))
		})
	})

	Context("Year / Month / Day / Weekday", func() {
		var v libtime.DateOrDateTime
		BeforeEach(func() {
			// 2026-01-15 is a Thursday
			v = libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
		})
		It("returns correct Year", func() {
			Expect(v.Year()).To(Equal(2026))
		})
		It("returns correct Month", func() {
			Expect(v.Month()).To(Equal(time.January))
		})
		It("returns correct Day", func() {
			Expect(v.Day()).To(Equal(15))
		})
		It("returns correct Weekday", func() {
			Expect(v.Weekday()).To(Equal(libtime.Weekday(time.Thursday)))
		})
	})

	Context("Hour / Minute / Second / Nanosecond", func() {
		Context("midnight UTC", func() {
			var v libtime.DateOrDateTime
			BeforeEach(func() {
				v = libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			})
			It("Hour returns 0", func() {
				Expect(v.Hour()).To(Equal(0))
			})
			It("Minute returns 0", func() {
				Expect(v.Minute()).To(Equal(0))
			})
			It("Second returns 0", func() {
				Expect(v.Second()).To(Equal(0))
			})
			It("Nanosecond returns 0", func() {
				Expect(v.Nanosecond()).To(Equal(0))
			})
		})
		Context("non-midnight", func() {
			var v libtime.DateOrDateTime
			BeforeEach(func() {
				v = libtime.DateOrDateTime(
					time.Date(2026, time.January, 15, 14, 30, 45, 123456789, time.UTC),
				)
			})
			It("Hour returns 14", func() {
				Expect(v.Hour()).To(Equal(14))
			})
			It("Minute returns 30", func() {
				Expect(v.Minute()).To(Equal(30))
			})
			It("Second returns 45", func() {
				Expect(v.Second()).To(Equal(45))
			})
			It("Nanosecond returns 123456789", func() {
				Expect(v.Nanosecond()).To(Equal(123456789))
			})
		})
	})

	Context("Time / TimePtr", func() {
		It("Time returns correct time.Time", func() {
			t := time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC)
			v := libtime.DateOrDateTime(t)
			Expect(v.Time().Equal(t)).To(BeTrue())
		})
		It("TimePtr on nil pointer returns nil", func() {
			var v *libtime.DateOrDateTime
			Expect(v.TimePtr()).To(BeNil())
		})
		It("TimePtr on non-nil pointer returns non-nil", func() {
			t := time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC)
			v := libtime.DateOrDateTime(t)
			Expect(v.TimePtr()).NotTo(BeNil())
			Expect(v.TimePtr().Equal(t)).To(BeTrue())
		})
	})

	Context("Format", func() {
		It("formats a midnight value as date-only", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			Expect(v.Format(time.DateOnly)).To(Equal("2026-01-15"))
		})
	})

	Context("Unix / UnixMicro", func() {
		var v libtime.DateOrDateTime
		BeforeEach(func() {
			v = libtime.DateOrDateTime(time.Unix(1687161394, 0))
		})
		It("Unix returns correct value", func() {
			Expect(v.Unix()).To(Equal(int64(1687161394)))
		})
		It("UnixMicro returns correct value", func() {
			Expect(v.UnixMicro()).To(Equal(int64(1687161394000000)))
		})
	})

	Context("Compare", func() {
		It("returns 0 for equal values", func() {
			a := libtime.DateOrDateTime(time.Unix(1000, 0))
			b := libtime.DateOrDateTime(time.Unix(1000, 0))
			Expect(a.Compare(b)).To(Equal(0))
		})
		It("returns -1 for earlier value", func() {
			a := libtime.DateOrDateTime(time.Unix(999, 0))
			b := libtime.DateOrDateTime(time.Unix(1000, 0))
			Expect(a.Compare(b)).To(Equal(-1))
		})
		It("returns 1 for later value", func() {
			a := libtime.DateOrDateTime(time.Unix(1000, 0))
			b := libtime.DateOrDateTime(time.Unix(999, 0))
			Expect(a.Compare(b)).To(Equal(1))
		})
	})

	Context("ComparePtr", func() {
		It("returns 0 when both nil", func() {
			var a *libtime.DateOrDateTime
			var b *libtime.DateOrDateTime
			Expect(a.ComparePtr(b)).To(Equal(0))
		})
		It("returns -1 when receiver is nil", func() {
			var a *libtime.DateOrDateTime
			b := libtime.DateOrDateTime(time.Unix(1000, 0)).Ptr()
			Expect(a.ComparePtr(b)).To(Equal(-1))
		})
		It("returns 1 when other is nil", func() {
			a := libtime.DateOrDateTime(time.Unix(1000, 0)).Ptr()
			var b *libtime.DateOrDateTime
			Expect(a.ComparePtr(b)).To(Equal(1))
		})
		It("returns 0 for equal non-nil values", func() {
			a := libtime.DateOrDateTime(time.Unix(1000, 0)).Ptr()
			b := libtime.DateOrDateTime(time.Unix(1000, 0)).Ptr()
			Expect(a.ComparePtr(b)).To(Equal(0))
		})
	})

	Context("Before / After / Equal / EqualPtr", func() {
		var earlier libtime.DateOrDateTime
		var later libtime.DateOrDateTime
		BeforeEach(func() {
			earlier = libtime.DateOrDateTime(
				time.Date(2026, time.January, 14, 0, 0, 0, 0, time.UTC),
			)
			later = libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
		})
		It("Before returns true for earlier value", func() {
			Expect(earlier.Before(later)).To(BeTrue())
		})
		It("Before returns false for later value", func() {
			Expect(later.Before(earlier)).To(BeFalse())
		})
		It("After returns true for later value", func() {
			Expect(later.After(earlier)).To(BeTrue())
		})
		It("After returns false for earlier value", func() {
			Expect(earlier.After(later)).To(BeFalse())
		})
		It("Equal returns true for same value", func() {
			a := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			Expect(later.Equal(a)).To(BeTrue())
		})
		It("Equal returns false for different value", func() {
			Expect(earlier.Equal(later)).To(BeFalse())
		})
		It("EqualPtr returns true when both nil", func() {
			var a *libtime.DateOrDateTime
			var b *libtime.DateOrDateTime
			Expect(a.EqualPtr(b)).To(BeTrue())
		})
		It("EqualPtr returns true for equal non-nil values", func() {
			a := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)).
				Ptr()
			b := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)).
				Ptr()
			Expect(a.EqualPtr(b)).To(BeTrue())
		})
		It("EqualPtr returns false when one is nil", func() {
			a := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)).
				Ptr()
			var b *libtime.DateOrDateTime
			Expect(a.EqualPtr(b)).To(BeFalse())
		})
	})

	Context("Add", func() {
		It("adding 1 hour to midnight value makes it non-midnight", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			result := v.Add(libtime.Duration(time.Hour))
			Expect(result.IsDateOnly()).To(BeFalse())
			Expect(result.Hour()).To(Equal(1))
		})
	})

	Context("Sub", func() {
		It("returns correct duration between two values", func() {
			a := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 2, 0, 0, 0, time.UTC))
			b := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			diff := a.Sub(b)
			Expect(diff.Duration()).To(Equal(2 * time.Hour))
		})
	})

	Context("AddDate", func() {
		var v libtime.DateOrDateTime
		BeforeEach(func() {
			v = libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
		})
		It("adds 1 year", func() {
			result := v.AddDate(1, 0, 0)
			Expect(result.Year()).To(Equal(2027))
		})
		It("adds 1 month", func() {
			result := v.AddDate(0, 1, 0)
			Expect(result.Month()).To(Equal(time.February))
		})
		It("adds 1 day", func() {
			result := v.AddDate(0, 0, 1)
			Expect(result.Day()).To(Equal(16))
		})
	})

	Context("Truncate", func() {
		It("truncates non-midnight to midnight UTC for that day", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 14, 30, 0, 0, time.UTC))
			result := v.Truncate(libtime.Duration(24 * time.Hour))
			Expect(result.IsDateOnly()).To(BeTrue())
			Expect(result.String()).To(Equal("2026-01-15"))
		})
	})

	Context("JSON struct round-trip", func() {
		type TestStruct struct {
			Date libtime.DateOrDateTime `json:"date"`
		}
		var original TestStruct
		var unmarshaled TestStruct
		var jsonBytes []byte
		BeforeEach(func() {
			original = TestStruct{
				Date: libtime.DateOrDateTime(
					time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC),
				),
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
		It("round-trips the value correctly", func() {
			Expect(unmarshaled.Date.String()).To(Equal(original.Date.String()))
			Expect(unmarshaled.Date.String()).To(Equal("2026-01-15"))
		})
	})

	Context("ParseDateOrDateTime", func() {
		Context("valid date string", func() {
			It("returns non-nil pointer, no error", func() {
				result, err := libtime.ParseDateOrDateTime(ctx, "2026-01-15")
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
			})
		})
		Context("valid RFC3339 string", func() {
			It("returns non-nil pointer, no error", func() {
				result, err := libtime.ParseDateOrDateTime(ctx, "2026-01-15T14:30:00Z")
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
			})
		})
		Context("invalid string", func() {
			It("returns nil, non-nil error", func() {
				result, err := libtime.ParseDateOrDateTime(ctx, "not-a-date")
				Expect(err).NotTo(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Context("ParseDateOrDateTimeDefault", func() {
		It("returns default for invalid input", func() {
			defaultVal := libtime.DateOrDateTime(
				time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			)
			result := libtime.ParseDateOrDateTimeDefault(ctx, "not-a-date", defaultVal)
			Expect(result).To(Equal(defaultVal))
		})
	})

	Context("DateOrDateTimePtr", func() {
		It("returns nil for nil input", func() {
			Expect(libtime.DateOrDateTimePtr(nil)).To(BeNil())
		})
		It("returns non-nil for non-nil input", func() {
			t := time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)
			result := libtime.DateOrDateTimePtr(&t)
			Expect(result).NotTo(BeNil())
		})
	})

	Context("DateOrDateTimes", func() {
		It("Interfaces returns correct slice", func() {
			values := libtime.DateOrDateTimes{
				libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)),
				libtime.DateOrDateTime(time.Date(2026, time.January, 16, 0, 0, 0, 0, time.UTC)),
			}
			ifaces := values.Interfaces()
			Expect(ifaces).To(HaveLen(2))
		})
		It("Strings returns correct slice", func() {
			values := libtime.DateOrDateTimes{
				libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)),
			}
			strs := values.Strings()
			Expect(strs).To(Equal([]string{"2026-01-15"}))
		})
	})

	Context("bytes.Buffer round-trip via MarshalText/UnmarshalText", func() {
		It("round-trips through encoding.TextMarshaler", func() {
			v := libtime.DateOrDateTime(time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC))
			b, err := v.MarshalText()
			Expect(err).To(BeNil())
			Expect(bytes.Equal(b, []byte("2026-01-15"))).To(BeTrue())

			var v2 libtime.DateOrDateTime
			err = v2.UnmarshalText(b)
			Expect(err).To(BeNil())
			Expect(v2.String()).To(Equal("2026-01-15"))
		})
	})
})
