// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"github.com/sofmon/wedit/model"
	"github.com/sofmon/wedit/renderer"
)

// Builder builds the wedit public folder
type Builder interface {
	ReadPageTemplate(path string) (string, error)
	ReadPageData(path string) (model.Page, error)
	WritePage(path string, page model.Page) error
	RebuildAll() error
}

type builder struct {
	cfg  Config
	rend renderer.Renderer
}

// NewBuilder creates new wedit Builder
func NewBuilder(cfg Config, rend renderer.Renderer) Builder {
	return &builder{cfg, rend}
}

/*
func (e *explorer) WriteAsset(path string, data []byte) error {

	dir, file := filepath.Split(path)

	folder := e.settings.Folders.Public + "/" + dir
	os.MkdirAll(folder, 0777)

	file = folder + "/" + file
	err := ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("Unable to save asset at '%v'. Error: %v", file, err)
		return err
	}

	return nil
}
*/
