// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"wedit/model"
)

// ReadPageData based on relative path
func ReadPageData(path string) (page model.Page, error error) {

	dir, file := filepath.Split(path)

	contentFolder := filepath.Join(cfg.ContentFolder, dir)
	contentFile := filepath.Join(contentFolder, file+".json")

	page = model.NewEmptyPage()

	data, err := ioutil.ReadFile(contentFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("✘ unable to load page data file '%v'. Error: %v\n", contentFile, err)
		return
	}

	if err == nil {
		err = json.Unmarshal(data, &page)
		if err != nil {
			log.Printf("✘ unable to load page data file '%v'. Error: %v\n", contentFile, err)
			return
		}
	}

	err = addRootData(&page)
	if err != nil {
		log.Printf("✘ unable to load root data. Error: %v\n", err)
		return
	}

	err = updateImagesSrcset(&page, path)
	if err != nil {
		log.Printf("✘ unable to load image data. Error: %v\n", err)
		return
	}

	return
}

func clearPublicDir(dir string) error {

	publicFolder := filepath.Join(cfg.PublicFolder, dir)
	os.MkdirAll(publicFolder, 0777)

	templateFolder, err := findTemplatePath(dir)
	if err != nil {
		return err
	}

	return copyDir(templateFolder, publicFolder, true)
}

// WritePage to the content and public folder
func WritePage(relPath string, page model.Page) error {

	log.Printf("✎ %s", relPath)

	dir, file := filepath.Split(relPath)

	err := clearPublicDir(dir)
	if err != nil {
		return err
	}

	publicFolder := filepath.Join(cfg.PublicFolder, dir)
	os.MkdirAll(publicFolder, 0777)

	contentFolder := filepath.Join(cfg.ContentFolder, dir)
	os.MkdirAll(contentFolder, 0777)

	contentFile := filepath.Join(contentFolder, file+".json")

	localData, rootData := splitRootData(page)

	data, err := json.MarshalIndent(localData, "", "  ")
	if err != nil {
		log.Printf("✘ unable to save page data at '%v'; error: %v", contentFile, err)
		return err
	}

	err = ioutil.WriteFile(contentFile, data, 0777)
	if err != nil {
		log.Printf("✘ unable to save page data at '%v'; error: %v", contentFile, err)
		return err
	}

	htmlFile := filepath.Join(publicFolder, file)
	templateFile, err := findTemplateFile(relPath)
	if err != nil {
		log.Printf("✘ unable to save page HTML at '%v'; error: %v", htmlFile, err)
		return err
	}

	err = renderHTML(htmlFile, templateFile, page)
	if err != nil {
		log.Printf("✘ unable to save page HTML at '%v'; error: %v", htmlFile, err)
		return err
	}

	wasUpdated, err := updateRootData(rootData)
	if err != nil {
		return err
	}

	if wasUpdated {
		log.Println("regenerating the whole public folder as root data was updated")
		err = RebuildAll()
		if err != nil {
			return err
		}
	} else {
		sm, err := loadSitemap()
		if err != nil {
			log.Printf("✘ unable to load sitemap due to an error: %v", err)
			return err
		}
		sm.Update(path.Join(cfg.SitemapHost, relPath))
	}

	return nil
}
