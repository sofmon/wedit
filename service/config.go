// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

// ShellCommand available to the user
type ShellCommand struct {
	Color  string   `json:"color"`
	Script []string `json:"script"`
}

// Config for the wedit http service
type Config struct {
	Host          string                  `json:"host"`
	Port          int                     `json:"port"`
	OpenBrowser   bool                    `json:"openBrowser"`
	EditAttr      string                  `json:"editAttr"`
	RepeatAttr    string                  `json:"repeatAttr"`
	MenuTextColor string                  `json:"menuTextColor"`
	ShellCommands map[string]ShellCommand `json:"shellCommands"`
}
