// Package main as wedit entry point
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config for the wedit (file) explorer
var cfg = struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	OpenBrowser bool   `json:"openBrowser"`
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
