// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	"sync"
	"time"
)

//counterfeiter:generate -o mocks/current-time-getter.go --fake-name CurrentTimeGetter . CurrentTimeGetter
type CurrentTimeGetter interface {
	Now() time.Time
}

type CurrentTimeGetterFunc func() DateTime

func (c CurrentTimeGetterFunc) Now() DateTime {
	return c()
}

//counterfeiter:generate -o mocks/current-time-setter.go --fake-name CurrentTimeSetter . CurrentTimeSetter
type CurrentTimeSetter interface {
	SetNow(now time.Time)
}

//counterfeiter:generate -o mocks/current-time.go --fake-name CurrentTime . CurrentTime
type CurrentTime interface {
	CurrentTimeGetter
	CurrentTimeSetter
}

func NewCurrentTime() CurrentTime {
	return &currentTime{}
}

type currentTime struct {
	mux sync.Mutex
	now *time.Time
}

func (n *currentTime) Now() time.Time {
	n.mux.Lock()
	defer n.mux.Unlock()
	if n.now != nil {
		return *n.now
	}
	return Now()
}

func (n *currentTime) SetNow(now time.Time) {
	n.mux.Lock()
	defer n.mux.Unlock()
	n.now = &now
}
