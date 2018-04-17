// Package explorer explores and manages the file system
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package explorer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sofmon/wedit/renderer"

	model "github.com/sofmon/wedit/model"
)

// Explorer read and write wedit related files
type Explorer interface {
	ReadPageTemplate(path string) (string, error)
	ReadPageData(path string) (model.Page, error)
	WritePageData(path string, page model.Page) error
}

type explorer struct {
	cfg  Config
	rend renderer.Renderer
}

func NewExplorer(cfg Config, rend renderer.Renderer) Explorer {
	return &explorer{cfg, rend}
}

const dataFileName = "index.json"
const templateFileName = "index.html"
const publicFileName = "index.html"

func (e *explorer) findTemplateFile(path string) string {
	subFolders := strings.Split(path, "/")
	folderToCheck := e.cfg.TemplateFolder
	templateFolder := e.cfg.TemplateFolder + "/"
	for i := 0; i < len(subFolders); i++ {
		folderToCheck += "/" + subFolders[i] + "/"
		stat, err := os.Stat(folderToCheck)
		if err == nil && stat.IsDir() {
			templateFolder = folderToCheck
		} else {
			break
		}
	}
	return templateFolder + templateFileName
}

func (e *explorer) ReadPageTemplate(path string) (string, error) {
	file := e.findTemplateFile(path)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Could not read page template file '%v'. Error: %v\n", file, err)
		return "", err
	}

	return string(data), nil
}

func (e *explorer) ReadPageData(path string) (page model.Page, error error) {

	file := e.cfg.PublicFolder + "/" + path + "/" + dataFileName
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

	return page, nil
}

func (e *explorer) WritePageData(path string, page model.Page) error {

	folder := e.cfg.PublicFolder + "/" + path
	os.MkdirAll(folder, 0777)

	file := folder + "/" + dataFileName

	data, err := json.Marshal(page)
	if err != nil {
		log.Printf("unable to save page data at '%v'; error: %v", file, err)
		return err
	}

	err = ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("unable to save page data at '%v'; error: %v", file, err)
		return err
	}

	indexFile := folder + "/" + publicFileName
	templateFile := e.findTemplateFile(path)

	err = e.rend.RenderHTML(indexFile, file, templateFile)
	if err != nil {
		log.Printf("unable to save page HTML at '%v'; error: %v", file, err)
		return err
	}

	return nil
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
