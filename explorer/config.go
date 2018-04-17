// Package explorer explores and manages the file system
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package explorer

// Config for the wedit (file) explorer
type Config struct {
	TemplateFolder string `json:"templateFolder"`
	PublicFolder   string `json:"publicFolder"`
}
