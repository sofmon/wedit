// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

// Config for the wedit (file) explorer
type Config struct {
	TemplateFolder   string `json:"templateFolder"`
	ContentFolder    string `json:"contentFolder"`
	PublicFolder     string `json:"publicFolder"`
	TemplateHTMLFile string `json:"templateHtmlFile"`
	PageHTMLFile     string `json:"pageHtmlFile"`
	PageJSONFile     string `json:"pageJsonFile"`
	RootJSONFile     string `json:"rootJsonFile"`
	RootKeyPrefix    string `json:"rootKeyPrefix"`
}
