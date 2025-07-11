// Copyright (c) 2024 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"time"

	libtime "github.com/bborbe/time"
)

func CreateTimeOfDay(startTime string, interval time.Duration, amount int) libtime.TimeOfDays {
	time := ParseTimeOfDay(startTime).Time(2000, 1, 1)
	result := []libtime.TimeOfDay{libtime.TimeOfDayFromTime(time)}
	for len(result) < amount {
		time = time.Add(interval)
		result = append(result, libtime.TimeOfDayFromTime(time))
	}
	return result
}
