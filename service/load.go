// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sofmon/wedit/builder"
	"github.com/sofmon/wedit/model"
)

func loadHandler(w http.ResponseWriter, r *http.Request) {

	path := getPathWithoutAction(r)

	page, err := builder.ReadPageData(path)
	if err != nil {
		log.Printf("unable to read page data path '%v'; error: %v", path, err)
		http.Error(w, "unable to read page data", http.StatusInternalServerError)
		return
	}

	pageWithSettings := model.PageWithSettings{
		Page: page,
		Settings: model.Settings{
			EditAttribute:   cfg.EditAttr,
			RepeatAttribute: cfg.RepeatAttr,
			MenuTextColor:   cfg.MenuTextColor,
		},
	}

	for k, v := range cfg.ShellCommands {
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
