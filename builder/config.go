// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

// Config for the wedit (file) explorer
type Config struct {
	TemplateFolder   string `json:"templateFolder"`
	ContentFolder    string `json:"contentFolder"`
	PublicFolder     string `json:"publicFolder"`
	PageFiles				 []string `json:"pageFiles"`
	RootJSONFile     string `json:"rootJsonFile"`
	RootKeyPrefix    string `json:"rootKeyPrefix"`
	EditAttr         string `json:"editAttribute"`
	RepeatAttr       string `json:"repeatAttribute"`
	IncludeAttr      string `json:"includeAttribute"`
	KeepWeditAttrs   bool   `json:"keepAttributes"`
}

func (c *Config) Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
