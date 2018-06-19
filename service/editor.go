// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"net/http"

	"github.com/sofmon/wedit/editor"
)

const (
	pageDataTemplate = `{"h":"","s":"","p":"","t":"","e":[],"r":[],"s":{"e":"","r":""}}`
)

func (s *service) editorHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Write([]byte(editor.EditorJSCode))
}
