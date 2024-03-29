// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typesinternal

import (
	"go/types"
)

func init() {
	SetGoVersion = func(conf *types.Config, version string) bool {
		conf.GoVersion = version
		return true
	}
}
