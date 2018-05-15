// Package config parses wedit.json config file
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sofmon/wedit/builder"
	"github.com/sofmon/wedit/service"
)

// Config for the wedit project
type Config struct {
	Builder builder.Config `json:"builder"`
	Service service.Config `json:"service"`
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

	if !strings.HasSuffix(cfg.Builder.ContentFolder, "/") {
		cfg.Builder.ContentFolder = cfg.Builder.ContentFolder + "/"
	}
	if !strings.HasSuffix(cfg.Builder.TemplateFolder, "/") {
		cfg.Builder.TemplateFolder = cfg.Builder.TemplateFolder + "/"
	}
	if !strings.HasSuffix(cfg.Builder.PublicFolder, "/") {
		cfg.Builder.PublicFolder = cfg.Builder.PublicFolder + "/"
	}

	cfg.Service.EditAttr = cfg.Builder.EditAttr
	cfg.Service.RepeatAttr = cfg.Builder.RepeatAttr

	return
}
