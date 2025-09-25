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

type TimeRange struct {
	From  stdtime.Time `json:"from,omitempty"`
	Until stdtime.Time `json:"until,omitempty"`
}

func (r TimeRange) Validate(ctx context.Context) error {
	return validation.All{
		validation.Name("range", validation.HasValidationFunc(func(ctx context.Context) error {
			if r.From.After(r.Until) {
				return errors.Wrapf(ctx, validation.Error, "from must be less than or equal to until")
			}
			return nil
		})),
	}.Validate(ctx)
}

func (r TimeRange) Ptr() *TimeRange {
	return &r
}
