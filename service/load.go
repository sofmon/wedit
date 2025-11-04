// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"encoding/json"
	"log"
	"net/http"

	"wedit/builder"
	"wedit/model"
)

func loadHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("p")

	_, path = verifyPath(path)

	page, err := builder.ReadPageData(path)
	if err != nil {
		log.Printf("✘ unable to read page data path '%v'; error: %v", path, err)
		http.Error(w, "unable to read page data", http.StatusInternalServerError)
		return
	}

	pageWithSettings := model.PageWithSettings{
		Page: page,
		Settings: model.Settings{
			EditAttribute:   cfg.EditAttr,
			RepeatAttribute: cfg.RepeatAttr,
			ClassAttribute:  cfg.ClassAttr,
			DarkMode:        cfg.DarkMode,
		},
	}

	data, err := json.Marshal(pageWithSettings)
	if err != nil {
		log.Printf("✘ unable to read page data path '%v'; error: %v", path, err)
		http.Error(w, "unable to read page data", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
