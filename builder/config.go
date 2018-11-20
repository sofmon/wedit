// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Config for the wedit (file) explorer
var cfg = struct {
	TemplateFolder      string   `json:"templateFolder"`
	ContentFolder       string   `json:"contentFolder"`
	PublicFolder        string   `json:"publicFolder"`
	VariationFolderName string   `json:"variationFolderName"`
	RootJSONFile        string   `json:"rootJsonFile"`
	RootKeyPrefix       string   `json:"rootKeyPrefix"`
	EditAttr            string   `json:"editAttribute"`
	RepeatAttr          string   `json:"repeatAttribute"`
	IncludeAttr         string   `json:"includeAttribute"`
	KeepWeditAttrs      bool     `json:"keepAttributes"`
	AllowedPageExt      []string `json:"allowedPageExt"`
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

	if !strings.HasSuffix(cfg.ContentFolder, "/") {
		cfg.ContentFolder = cfg.ContentFolder + "/"
	}
	if !strings.HasSuffix(cfg.TemplateFolder, "/") {
		cfg.TemplateFolder = cfg.TemplateFolder + "/"
	}
	if !strings.HasSuffix(cfg.PublicFolder, "/") {
		cfg.PublicFolder = cfg.PublicFolder + "/"
	}

	return nil
}
