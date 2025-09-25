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

// DateRangeFromTime creates a DateRange from two time.Time values.
// It converts the from and until times to Date types and returns a DateRange.
func DateRangeFromTime(from, until stdtime.Time) DateRange {
	return DateRange{
		From:  Date(from),
		Until: Date(until),
	}
}

type DateRange struct {
	From  Date `json:"from,omitempty"`
	Until Date `json:"until,omitempty"`
}

func (r DateRange) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("from", r.From),
		validation.Name("until", r.Until),
		validation.Name("range", validation.HasValidationFunc(func(ctx context.Context) error {
			if r.From.Time().After(r.Until.Time()) {
				return errors.Wrapf(ctx, validation.Error, "from must be less than or equal to until")
			}
			return nil
		})),
	}.Validate(ctx)
}

func (r DateRange) Ptr() *DateRange {
	return &r
}
