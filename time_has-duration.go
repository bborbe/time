// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import stdtime "time"

//counterfeiter:generate -o mocks/has-duration.go --fake-name HasDuration . HasDuration
type HasDuration interface {
	Duration() stdtime.Duration
}
