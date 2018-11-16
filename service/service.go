// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"fmt"
	"net/http"
)

// Service processing wedit HTTP requests
type Service interface {
	ListenAndServe() error
}

var (
	staticHandler http.Handler
)

// ListenAndServe blocks the current go routine and start serving the wedit HTTP request
func ListenAndServe() error {

	staticHandler = http.FileServer(http.Dir(cfg.PublicFolder))

	/*
		By default ActionPAth os '~', so:
		GET  ~.js			- editor js
		GET  ~?p=/some/url	- load page data
		PUT	 ~?p=/some/url	- save page data
		POST ~?p=/some/url	- upload asset (image)
		GET	 /				- Serve page or asset
	*/

	http.HandleFunc(
		"/~.js",
		editorHandler, // editor.go
	)

	http.HandleFunc(
		"/~",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				loadHandler(w, r) // load.go
				return
			case http.MethodPut:
				saveHandler(w, r) // save.go
				return
			case http.MethodPost:
				imageHandler(w, r) // image.go
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
		},
	)

	http.HandleFunc("/", pageHandler) // page.go

	endPoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return http.ListenAndServe(endPoint, nil)
}
