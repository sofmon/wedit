// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config for the wedit http service
var cfg = struct {
	PublicFolder   string   `json:"publicFolder"`
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	OpenBrowser    bool     `json:"openBrowser"`
	EditAttr       string   `json:"editAttribute"`
	RepeatAttr     string   `json:"repeatAttribute"`
	DefaultPage    string   `json:"defaultPage"`
	AllowedPageExt []string `json:"allowedPageExt"`
}{}

// LoadConfig form wedit.json file
func LoadConfig(file string) error {

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
