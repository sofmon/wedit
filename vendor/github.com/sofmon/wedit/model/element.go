// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

// Element represents an html element on a page
type Element struct {
	Key  Key    `json:"k"`
	Text string `json:"t"`
}
