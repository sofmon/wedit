// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

// Command available for the user
type Command struct {
	Name  string `json:"n"`
	Color string `json:"c"`
}
