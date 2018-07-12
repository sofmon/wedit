// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

// Image represents an image on a page
type Image struct {
	Key Key    `json:"k"`
	Src string `json:"s"`
}
