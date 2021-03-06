// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

// Repeat represents a repeatable element on a page
type Repeat struct {
	Key      Key   `json:"k"`
	CopyKeys []Key `json:"c"`
}
