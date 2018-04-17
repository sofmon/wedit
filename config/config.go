// Package config parses wedit.json config file
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/sofmon/wedit/explorer"
	"github.com/sofmon/wedit/renderer"
	"github.com/sofmon/wedit/service"
)

// Config for the wedit project
type Config struct {
	Explorer explorer.Config `json:"explorer"`
	Service  service.Config  `json:"service"`
	Renderer renderer.Config `json:"renderer"`
	Edit     struct {
		OpenBrowser bool `json:"openBrowser"`
	} `json:"edit"`
}

const configFile = "wedit.json"

// LoadConfig from the default wedit.json file
func LoadConfig() (cfg Config, err error) {

	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		err = fmt.Errorf("config: unable to read config file (%s) due to error: %v", configFile, err)
		return
	}

	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		err = fmt.Errorf("config: unable to parse config file (%s) as JSON due to error: %v", configFile, err)
		return
	}

	return
}
