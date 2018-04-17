// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sofmon/wedit/model"
)

func (s *service) saveHandler(w http.ResponseWriter, r *http.Request) {

	path := getPathWithoutAction(r)

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

	err = s.ex.WritePageData(path, page)
	if err != nil {
		log.Printf("unable to save page for path '%v'; error: %v", path, err)
		http.Error(w, "unable to save page", http.StatusInternalServerError)
		return
	}
}
