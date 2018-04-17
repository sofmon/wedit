// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const editorScriptTag = `<script async="" src="/!editor.js"></script>`

func (s *service) pageHandler(w http.ResponseWriter, r *http.Request) {

	path := getPathWithoutAction(r)

	if filepath.Ext(path) != "" {
		r.URL.RawPath = strings.Replace(r.URL.RawPath, "/!/", "/", 1)
		r.URL.Path = strings.Replace(r.URL.Path, "/!/", "/", 1)
		s.staticHandler.ServeHTTP(w, r)
		return
	}

	if !strings.HasSuffix(r.URL.RawPath, "/") {
		http.Redirect(w, r, r.URL.RawPath+"/", http.StatusPermanentRedirect)
		return
	}

	templateRaw, err := s.ex.ReadPageTemplate(path)
	if err != nil {
		log.Printf("unable to serve template on path '%v'. Error: %v", path, err)
		http.NotFound(w, r)
		return
	}

	templateSplit := strings.Split(templateRaw, "</body>")
	templateWithEditor := templateSplit[0] + editorScriptTag + "</body>"

	w.Write([]byte(templateWithEditor))
}
