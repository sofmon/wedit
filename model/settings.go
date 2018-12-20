// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

// Settings for the wedit page
type Settings struct {
	EditAttribute   string    `json:"e"`
	RepeatAttribute string    `json:"r"`
	MenuTextColor   string    `json:"m"`
	Commands        []Command `json:"c"`
	DarkMode        bool      `json:"d"`
}
