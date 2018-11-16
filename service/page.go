// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sofmon/wedit/builder"
)

const editorScriptTag = `<script async="" src="/~.js"></script>`

func pageHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	isHTML, path := verifyPath(path)

	if !isHTML {
		staticHandler.ServeHTTP(w, r)
		return
	}

	templateHTML, err := builder.ReadPageTemplate(path)
	if err != nil {
		log.Printf("unable to serve template on path '%v'. Error: %v", path, err)
		http.NotFound(w, r)
		return
	}

	templateSplit := strings.Split(templateHTML, "</body>")
	templateWithEditor := templateSplit[0] + editorScriptTag + "</body>"

	w.Write([]byte(templateWithEditor))
}

func verifyPath(path string) (isHTML bool, exactPath string) {

	if strings.HasSuffix(path, "/") {
		return true, path + cfg.DefaultPage
	}

	ext := filepath.Ext(path)
	for _, e := range cfg.AllowedPageExt {
		if e == ext {
			return true, path
		}
	}

	return false, path
}
