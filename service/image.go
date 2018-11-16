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

	"github.com/sofmon/wedit/builder"
	"github.com/sofmon/wedit/model"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("p")

	path = filepath.Dir(path)

	imageName := r.URL.Query().Get("n")
	imageKey := model.Key(r.URL.Query().Get("k"))

	imageData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("unable to process image ('%s') data; error: %v", imageName, err)
		http.Error(w, "unable to process image", http.StatusInternalServerError)
		return
	}

	img, err := builder.WriteImage(path, imageKey, imageName, imageData)
	if err != nil {
		log.Printf("unable to process image ('%s') data; error: %v", imageName, err)
		http.Error(w, "unable to process image", http.StatusInternalServerError)
		return
	}

	resBody, err := json.Marshal(img)
	if err != nil {
		log.Printf("unable to process image ('%s') data; error: %v", imageName, err)
		http.Error(w, "unable to process image", http.StatusInternalServerError)
		return
	}

	w.Write(resBody)
}
