// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sofmon/wedit/model"
)

func (s *service) saveHandler(w http.ResponseWriter, r *http.Request) {

	path := getPathWithoutAction(r)
	pageFile := "index"
	builderConfig := s.bld.GetConfig()
	if filepath.Ext(path) != "" {
		pageFile = strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))
	  if !builderConfig.Contains(builderConfig.PageFiles, pageFile) {
			pageFile = "index"
		}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("unable to read save body for path '%v'; error: %v", path, err)
		http.Error(w, "unable to read save body", http.StatusInternalServerError)
		return
	}

	var page model.Page
	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Printf("unable to marshal save body for path '%v'; error: %v", path, err)
		http.Error(w, "unable to marshal save body", http.StatusInternalServerError)
		return
	}

	dir := filepath.Dir(path)
	if dir == "." {
		dir = ""
	} else {
		dir += "/"
	}
	oldPage, err := s.bld.ReadPageData(dir, pageFile)
	if err != nil {
		log.Printf("unable to read old page data path '%v'; error: %v", path, err)
		http.Error(w, "unable to read old page data ", http.StatusInternalServerError)
		return
	}

	for k, v := range oldPage.Elements {
		if _, ok := page.Elements[k]; !ok {
			page.Elements[k] = v
		}
	}

	for k, v := range oldPage.Images {
		if _, ok := page.Images[k]; !ok {
			page.Images[k] = v
		}
	}

	for k, v := range oldPage.Repeats {
		if _, ok := page.Repeats[k]; !ok {
			page.Repeats[k] = v
		}
	}

	err = s.bld.WritePage(dir, page, pageFile)
	if err != nil {
		log.Printf("unable to save page for path '%v'; error: %v", path, err)
		http.Error(w, "unable to save page", http.StatusInternalServerError)
		return
	}
}
