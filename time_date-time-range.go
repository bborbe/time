// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"context"

	"github.com/bborbe/errors"
	"github.com/bborbe/validation"
)

type DateTimeRange struct {
	From  DateTime `json:"from,omitempty"`
	Until DateTime `json:"until,omitempty"`
}

func (r DateTimeRange) Validate(ctx context.Context) error {
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

func (r DateTimeRange) Ptr() *DateTimeRange {
	return &r
}
