// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"context"
	"encoding"
	"encoding/json"
	"strings"
	stdtime "time"

	"github.com/bborbe/errors"
	"github.com/bborbe/parse"
	"github.com/bborbe/validation"
)

type DateOrDateTimes []DateOrDateTime

func (d DateOrDateTimes) Interfaces() []interface{} {
	result := make([]interface{}, len(d))
	for i, ss := range d {
		result[i] = ss
	}
	return result
}

func (d DateOrDateTimes) Strings() []string {
	result := make([]string, len(d))
	for i, ss := range d {
		result[i] = ss.String()
	}
	return result
}

func DateOrDateTimeFromBinary(ctx context.Context, value []byte) (*DateOrDateTime, error) {
	var t stdtime.Time
	if err := t.UnmarshalBinary(value); err != nil {
		return nil, errors.Wrapf(ctx, err, "unmarshalBinary failed")
	}
	return DateOrDateTime(t).Ptr(), nil
}

func ParseDateOrDateTimeDefault(
	ctx context.Context,
	value interface{},
	defaultValue DateOrDateTime,
) DateOrDateTime {
	result, err := ParseDateOrDateTime(ctx, value)
	if err != nil {
		return defaultValue
	}
	return *result
}

func ParseDateOrDateTime(ctx context.Context, value interface{}) (*DateOrDateTime, error) {
	str, err := parse.ParseString(ctx, value)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "parse value failed")
	}
	t, err := ParseTime(ctx, str)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "parse time failed")
	}
	return DateOrDateTimePtr(t), nil
}

func DateOrDateTimePtr(value *stdtime.Time) *DateOrDateTime {
	if value == nil {
		return nil
	}
	return DateOrDateTime(*value).Ptr()
}

// NewDateOrDateTime creates a DateOrDateTime representing the date and time specified by the given parameters.
// It wraps the standard library's time.Date function with the same parameter signature.
func NewDateOrDateTime(
	year int,
	month stdtime.Month,
	day, hour, min, sec, nsec int,
	loc *stdtime.Location,
) DateOrDateTime {
	return DateOrDateTime(stdtime.Date(year, month, day, hour, min, sec, nsec, loc))
}

type DateOrDateTime stdtime.Time

var _ encoding.TextMarshaler = DateOrDateTime{}

var _ encoding.TextUnmarshaler = (*DateOrDateTime)(nil)

// isMidnightUTC reports whether t is exactly midnight UTC (all time components zero in UTC).
// This is the key discriminator for the round-trip serialization rule.
func isMidnightUTC(t stdtime.Time) bool {
	u := t.UTC()
	return u.Hour() == 0 && u.Minute() == 0 && u.Second() == 0 && u.Nanosecond() == 0
}

func (d DateOrDateTime) Year() int {
	return d.Time().Year()
}

func (d DateOrDateTime) Month() stdtime.Month {
	return d.Time().Month()
}

func (d DateOrDateTime) Day() int {
	return d.Time().Day()
}

func (d DateOrDateTime) Weekday() Weekday {
	return Weekday(d.Time().Weekday())
}

func (d DateOrDateTime) Hour() int {
	return d.Time().Hour()
}

func (d DateOrDateTime) Minute() int {
	return d.Time().Minute()
}

func (d DateOrDateTime) Second() int {
	return d.Time().Second()
}

func (d DateOrDateTime) Nanosecond() int {
	return d.Time().Nanosecond()
}

func (d DateOrDateTime) String() string {
	t := d.Time()
	if t.IsZero() {
		return ""
	}
	if isMidnightUTC(t) {
		return t.UTC().Format(stdtime.DateOnly)
	}
	return t.Format(stdtime.RFC3339Nano)
}

func (d DateOrDateTime) Validate(ctx context.Context) error {
	if d.Time().IsZero() {
		return errors.Wrapf(ctx, validation.Error, "time is zero")
	}
	return nil
}

func (d DateOrDateTime) Ptr() *DateOrDateTime {
	return &d
}

// IsZero reports whether d represents the zero time instant.
func (d DateOrDateTime) IsZero() bool {
	return d.Time().IsZero()
}

func (d DateOrDateTime) UTC() DateOrDateTime {
	return DateOrDateTime(d.Time().UTC())
}

func (d DateOrDateTime) Clone() DateOrDateTime {
	return d
}

func (d *DateOrDateTime) ClonePtr() *DateOrDateTime {
	if d == nil {
		return nil
	}
	return d.Clone().Ptr()
}

func (d *DateOrDateTime) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	switch str {
	case "", "null":
		*d = DateOrDateTime(stdtime.Time{})
		return nil
	default:
		t, err := ParseTime(context.Background(), str)
		if err != nil {
			return errors.Wrapf(context.Background(), err, "parse time failed")
		}
		*d = DateOrDateTime(*t)
		return nil
	}
}

func (d DateOrDateTime) MarshalJSON() ([]byte, error) {
	t := d.Time()
	if t.IsZero() {
		return json.Marshal(nil)
	}
	if isMidnightUTC(t) {
		return json.Marshal(t.UTC().Format(stdtime.DateOnly))
	}
	return json.Marshal(t.Format(stdtime.RFC3339Nano))
}

func (d DateOrDateTime) MarshalText() ([]byte, error) {
	t := d.Time()
	if t.IsZero() {
		return nil, nil
	}
	if isMidnightUTC(t) {
		return []byte(t.UTC().Format(stdtime.DateOnly)), nil
	}
	return []byte(t.Format(stdtime.RFC3339Nano)), nil
}

func (d *DateOrDateTime) UnmarshalText(b []byte) error {
	str := string(b)
	if len(str) == 0 {
		*d = DateOrDateTime(stdtime.Time{})
		return nil
	}
	t, err := ParseTime(context.Background(), str)
	if err != nil {
		return errors.Wrapf(context.Background(), err, "parse time failed")
	}
	*d = DateOrDateTime(*t)
	return nil
}

func (d DateOrDateTime) MarshalBinary() ([]byte, error) {
	return d.Time().MarshalBinary()
}

func (d DateOrDateTime) Time() stdtime.Time {
	return stdtime.Time(d)
}

func (d *DateOrDateTime) TimePtr() *stdtime.Time {
	if d == nil {
		return nil
	}
	t := stdtime.Time(*d)
	return &t
}

func (d DateOrDateTime) Format(layout string) string {
	return d.Time().Format(layout)
}

func (d DateOrDateTime) Unix() int64 {
	return d.Time().Unix()
}

func (d DateOrDateTime) UnixMicro() int64 {
	return d.Time().UnixMicro()
}

func (d DateOrDateTime) Compare(other DateOrDateTime) int {
	return Compare(d.Time(), other.Time())
}

func (d *DateOrDateTime) ComparePtr(other *DateOrDateTime) int {
	if d == nil && other == nil {
		return 0
	}
	if d == nil {
		return -1
	}
	if other == nil {
		return 1
	}
	return d.Compare(*other)
}

func (d DateOrDateTime) Before(other HasTime) bool {
	return d.Time().Before(other.Time())
}

func (d DateOrDateTime) After(other HasTime) bool {
	return d.Time().After(other.Time())
}

func (d DateOrDateTime) Equal(other DateOrDateTime) bool {
	return d.Time().Equal(other.Time())
}

func (d *DateOrDateTime) EqualPtr(other *DateOrDateTime) bool {
	if d == nil && other == nil {
		return true
	}
	if d != nil && other != nil {
		return d.Equal(*other)
	}
	return false
}

func (d DateOrDateTime) Add(duration HasDuration) DateOrDateTime {
	return DateOrDateTime(d.Time().Add(duration.Duration()))
}

func (d DateOrDateTime) Sub(other HasTime) Duration {
	return Duration(d.Time().Sub(other.Time()))
}

func (d DateOrDateTime) AddDate(years int, months int, days int) DateOrDateTime {
	return DateOrDateTime(d.Time().AddDate(years, months, days))
}

// Deprecated: Use AddDate instead.
// AddTime adds the given years, months, and days to the DateOrDateTime but will be removed in future versions.
func (d DateOrDateTime) AddTime(years int, months int, days int) DateOrDateTime {
	return d.AddDate(years, months, days)
}

func (d DateOrDateTime) Truncate(duration HasDuration) DateOrDateTime {
	return DateOrDateTime(d.Time().Truncate(duration.Duration()))
}

// IsDateOnly reports whether d represents a date-only value (midnight UTC).
func (d DateOrDateTime) IsDateOnly() bool {
	return !d.IsZero() && isMidnightUTC(d.Time())
}

// AsDate returns the date component of d as a Date value.
func (d DateOrDateTime) AsDate() Date {
	return ToDate(d.Time())
}

// AsDateTime returns d as a DateTime value.
func (d DateOrDateTime) AsDateTime() DateTime {
	return DateTime(d.Time())
}
