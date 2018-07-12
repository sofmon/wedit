// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

import (
	"strings"
)

// Key for wedit element
type Key string

// IsGlobal indicate if key is global for the whole site
func (k Key) IsGlobal() bool {
	return strings.HasPrefix(string(k), "!")
}
