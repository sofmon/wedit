// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

// Class represents a CSS class selection for an element
type Class struct {
	Key   Key    `json:"k"`
	Value string `json:"v"`
}
