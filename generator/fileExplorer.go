// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const dataFileName = "index.json"
const templateFileName = "index.html"
const publicFileName = "index.html"

// FileExplorer is used for i/o operations on the file system
type FileExplorer struct {
	settings Settings
}

// NewFileExplorer creates new file explorer
func NewFileExplorer(settings Settings) FileExplorer {
	return FileExplorer{settings}
}

func (f FileExplorer) findTemplateFile(path string) string {
	subFolders := strings.Split(path, "/")
	folderToCheck := f.settings.Folders.Template
	templateFolder := f.settings.Folders.Template + "/"
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

func (f FileExplorer) ReadPageTemplate(path string) (string, error) {
	file := f.findTemplateFile(path)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Could not read page template file '%v'. Error: %v\n", file, err)
		return "", err
	}

	return string(data), nil
}

func (f FileExplorer) ReadPageData(path string) (Page, error) {

	file := f.settings.Folders.Public + "/" + path + "/" + dataFileName
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return Page{}, nil
		}
		log.Printf("Could not load page data file '%v'. Error: %v\n", file, err)
		return Page{}, err
	}

	var page Page
	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("Could not load page data file '%v'. Error: %v\n", file, err)
		return Page{}, err
	}

	return page, nil
}

func (f FileExplorer) WritePageData(path string, page Page) error {

	folder := f.settings.Folders.Public + "/" + path
	os.MkdirAll(folder, 0777)

	file := folder + "/" + dataFileName

	data, err := json.Marshal(page)
	if err != nil {
		log.Printf("Unable to save page data at '%v'. Error: %v", file, err)
		return err
	}

	err = ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("Unable to save page data at '%v'. Error: %v", file, err)
		return err
	}

	indexFile := folder + "/" + publicFileName
	templateFile := f.findTemplateFile(path)

	err = writePageAsHtml(indexFile, templateFile, &page)
	//err = ioutil.WriteFile(file, []byte(page.HTML), 0777)
	if err != nil {
		log.Printf("Unable to save page HTML at '%v'. Error: %v", file, err)
		return err
	}

	return nil
}

func (f FileExplorer) WriteAsset(path string, data []byte) error {

	dir, file := filepath.Split(path)

	folder := f.settings.Folders.Public + "/" + dir
	os.MkdirAll(folder, 0777)

	file = folder + "/" + file
	err := ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("Unable to save asset at '%v'. Error: %v", file, err)
		return err
	}

	return nil
}
