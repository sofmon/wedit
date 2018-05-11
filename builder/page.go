// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/sofmon/wedit/model"
)

func (b *builder) ReadPageData(path string) (page model.Page, error error) {

	file := b.cfg.ContentFolder + "/" + path + "/" + b.cfg.PageJSONFile
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			page = model.Page{}
			err = nil // it is an empty data
			return
		}
		log.Printf("Could not load page data file '%v'. Error: %v\n", file, err)
		return
	}

	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("Could not load page data file '%v'. Error: %v\n", file, err)
		return
	}

	return b.addRootData(page)
}

func (b *builder) clearPublic(path string) error {

	publicFolder := b.cfg.PublicFolder + "/" + path
	os.MkdirAll(publicFolder, 0777)

	templateFolder := b.findTemplatePath(path)

	return copyDir(templateFolder, publicFolder)
}

func (b *builder) WritePage(path string, page model.Page) error {

	err := b.clearPublic(path)
	if err != nil {
		return err
	}

	publicFolder := b.cfg.PublicFolder + "/" + path
	os.MkdirAll(publicFolder, 0777)

	contentFolder := b.cfg.ContentFolder + "/" + path
	os.MkdirAll(contentFolder, 0777)

	file := contentFolder + "/" + b.cfg.PageJSONFile

	localData, rootData := b.splitRootData(page)

	data, err := json.MarshalIndent(localData, "", "  ")
	if err != nil {
		log.Printf("unable to save page data at '%v'; error: %v", file, err)
		return err
	}

	err = ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("unable to save page data at '%v'; error: %v", file, err)
		return err
	}

	indexFile := publicFolder + "/" + b.cfg.PageHTMLFile
	templateFile := b.findTemplateFile(path)

	err = b.renderHTML(indexFile, templateFile, page)
	if err != nil {
		log.Printf("unable to save page HTML at '%v'; error: %v", file, err)
		return err
	}

	wasUpdated, err := b.updateRootData(rootData)
	if err != nil {
		return err
	}

	if wasUpdated {
		log.Println("regenerating the whole public folder as root data was updated")
		err = b.RebuildAll()
		if err != nil {
			return err
		}
	}

	return nil
}
