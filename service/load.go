// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"github.com/sofmon/wedit/model"
	"strings"
)

func (s *service) loadHandler(w http.ResponseWriter, r *http.Request) {

	path := getPathWithoutAction(r)
	pageFile := "index"
	builderConfig := s.bld.GetConfig()
	if filepath.Ext(path) != "" {
		pageFile = strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))
	  if !builderConfig.Contains(builderConfig.PageFiles, pageFile) {
			pageFile = "index"
		}
	}

	dir := filepath.Dir(path)
	if dir == "." {
		dir = ""
	} else {
		dir += "/"
	}
	page, err := s.bld.ReadPageData(dir, pageFile)
	if err != nil {
		log.Printf("unable to read page data path '%v'; error: %v", path, err)
		http.Error(w, "unable to read page data", http.StatusInternalServerError)
		return
	}

	pageWithSettings := model.PageWithSettings{
		Page: page,
		Settings: model.Settings{
			EditAttribute:   s.cfg.EditAttr,
			RepeatAttribute: s.cfg.RepeatAttr,
			MenuTextColor:   s.cfg.MenuTextColor,
		},
	}

	for k, v := range s.cfg.ShellCommands {
		pageWithSettings.Settings.Commands = append(
			pageWithSettings.Settings.Commands,
			model.Command{
				Name:  k,
				Color: v.Color,
			},
		)
	}

	data, err := json.Marshal(pageWithSettings)
	if err != nil {
		log.Printf("unable to read page data path '%v'; error: %v", path, err)
		http.Error(w, "unable to read page data", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
