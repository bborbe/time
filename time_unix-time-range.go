// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"context"
	stdtime "time"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

// UnixTimeRangeFromTime creates a UnixTimeRange from two time.Time values.
// It converts the from and until times to UnixTime types and returns a UnixTimeRange.
func UnixTimeRangeFromTime(from, until stdtime.Time) UnixTimeRange {
	return UnixTimeRange{
		From:  UnixTime(from),
		Until: UnixTime(until),
	}
}

type UnixTimeRange struct {
	From  UnixTime `json:"from,omitempty"`
	Until UnixTime `json:"until,omitempty"`
}

func (r UnixTimeRange) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("from", r.From),
		validation.Name("until", r.Until),
		validation.Name("range", validation.HasValidationFunc(func(ctx context.Context) error {
			if r.From.After(r.Until) {
				return errors.Wrapf(ctx, validation.Error, "from must be less than or equal to until")
			}
			return nil
		})),
	}.Validate(ctx)
}

func (r UnixTimeRange) Ptr() *UnixTimeRange {
	return &r
}
