// Package renderer generates HTML files in the `public` folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package renderer

// Config for the wedit renderer
type Config struct {
	EditAttr   string `json:"editAttribute"`
	RepeatAttr string `json:"repeatAttribute"`
	KeepAttr   bool   `json:"keepAttribute"`
}
