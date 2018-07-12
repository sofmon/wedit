// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"github.com/sofmon/wedit/model"
)

// Builder builds the wedit public folder
type Builder interface {
	ReadPageTemplate(path string) (string, error)
	ReadPageData(path string) (model.Page, error)
	WritePage(path string, page model.Page) error
	WriteImage(path string, key model.Key, name string, data []byte) (srcset []string, err error)
	RebuildAll() error
}

type builder struct {
	cfg Config
}

// NewBuilder creates new wedit Builder
func NewBuilder(cfg Config) Builder {
	return &builder{cfg}
}
