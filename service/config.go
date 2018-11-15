// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ShellCommand available to the user
type ShellCommand struct {
	Color  string   `json:"color"`
	Script []string `json:"script"`
}

// Config for the wedit http service
var cfg = struct {
	PublicFolder  string                  `json:"publicFolder"`
	Cookie        string                  `json:"cookie"`
	Host          string                  `json:"host"`
	Port          int                     `json:"port"`
	OpenBrowser   bool                    `json:"openBrowser"`
	EditAttr      string                  `json:"editAttr"`
	RepeatAttr    string                  `json:"repeatAttr"`
	MenuTextColor string                  `json:"menuTextColor"`
	ShellCommands map[string]ShellCommand `json:"shellCommands"`
}{}

// LoadConfig form wedit.json file
func LoadConfig() error {

	file := "wedit.json"

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("config: unable to read config file (%s) due to error: %v", file, err)
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("config: unable to parse config file (%s) as JSON due to error: %v", file, err)
	}

	return nil
}
