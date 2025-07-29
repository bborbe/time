// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"context"
	"strconv"
	"time"

	"github.com/bborbe/errors"
	libparse "github.com/bborbe/parse"
)

const (
	DateLayout   Layout = "2006-01-02"
	SecondLayout Layout = "second"
	MilliLayout  Layout = "milli"
	MicroLayout  Layout = "micro"
	NanoLayout   Layout = "nano"
	RFC3339      Layout = time.RFC3339
	RFC3339Nano  Layout = time.RFC3339Nano
)

type Layouts []Layout

func (l Layouts) Parse(ctx context.Context, value interface{}) (*time.Time, error) {
	for _, ll := range l {
		if v, err := ll.Parse(ctx, value); err == nil {
			return v, nil
		}
	}
	return nil, errors.Errorf(ctx, "parse '%v' with any layouts failed", value)
}

// Layout is one of (millis,seconds,nano) or any time.Layout like time.RFC3339Nano
type Layout string

func (l Layout) String() string {
	return string(l)
}

func (l Layout) Format(time time.Time) string {
	switch l {
	case SecondLayout:
		return strconv.FormatInt(time.Unix(), 10)
	case MilliLayout:
		return strconv.FormatInt(time.UnixMilli(), 10)
	case MicroLayout:
		return strconv.FormatInt(time.UnixMicro(), 10)
	case NanoLayout:
		return strconv.FormatInt(time.UnixNano(), 10)
	default:
		return time.Format(l.String())
	}
}

func (l Layout) Parse(ctx context.Context, value interface{}) (*time.Time, error) {
	switch l {
	case SecondLayout:
		i, err := libparse.ParseInt64(ctx, value)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "convert to int failed")
		}
		t := time.Unix(i, 0)
		return &t, nil
	case MilliLayout:
		i, err := libparse.ParseInt64(ctx, value)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "convert to int failed")
		}
		t := time.UnixMilli(i)
		return &t, nil
	case MicroLayout:
		i, err := libparse.ParseInt64(ctx, value)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "convert to int failed")
		}
		t := time.UnixMicro(i)
		return &t, nil
	case NanoLayout:
		i, err := libparse.ParseInt64(ctx, value)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "convert to int failed")
		}
		t := time.Unix(i/int64(time.Second), i%int64(time.Second))
		return &t, nil
	default:
		switch v := value.(type) {
		case string:
			t, err := time.Parse(l.String(), v)
			if err != nil {
				return nil, errors.Wrapf(ctx, err, "parse '%s' with layout '%s' failed", v, l)
			}
			return &t, nil
		default:
			return nil, errors.Errorf(ctx, "can not parse %T with layout '%s' failed", value, l)
		}
	}
}
