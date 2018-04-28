// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sofmon/wedit/builder"
)

// Service processing wedit HTTP requests
type Service interface {
	ListenAndServe() error
}

type service struct {
	cfg           Config
	pub           string
	bld           builder.Builder
	staticHandler http.Handler
}

// NewService for processing wedit HTTP requests
func NewService(cfg Config, pub string, bld builder.Builder) Service {
	return &service{cfg, pub, bld, http.FileServer(http.Dir(pub))}
}

// ListenAndServe blocks the current go routine and start serving the wedit HTTP request
func (s *service) ListenAndServe() error {

	http.HandleFunc("/!/", s.pageHandler)           // page.go
	http.HandleFunc("/!editor.js", s.editorHandler) // editor.go
	http.HandleFunc("/!save/", s.saveHandler)       // save.go
	http.HandleFunc("/!picture/", s.pictureHandler) // picture.go
	http.HandleFunc("/!custom/", s.customHandler)   // custom.go

	http.Handle("/", s.staticHandler)

	endPoint := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	return http.ListenAndServe(endPoint, nil)
}

// Action to be performed by the current request
type Action int

const (
	// ActionGetPage is used to retrieve page HTML
	ActionGetPage Action = iota
	// ActionServeResource is used when HTML or static resource is requested
	ActionServeResource
	// ActionJsRequest is used when the editor JS file is requested
	ActionJsRequest
	// ActionSavePage is used when page is being saved by the editor
	ActionSavePage
	// ActionLoadPage is used when page is being loaded by the editor
	ActionLoadPage
	// ActionUploadImage is used when image is uploaded by the editor
	ActionUploadImage
)

const (
	actionJsRequestPathKey   = "!editor.js"
	actionEditPagePathKey    = "!"
	actionSavePagePathKey    = "!save"
	actionLoadPagePathKey    = "!load"
	actionUploadImagePathKey = "!upload"
)

// getPathWithoutAction retrieves the path of the requested resource excluding the wedit parameters/action
func getPathWithoutAction(r *http.Request) string {

	a := getAction(r)

	switch a {

	case ActionJsRequest:
		ref, err := url.Parse(r.Referer())
		if err != nil {
			return "/" // TODO: evaluate
		}
		return strings.Replace(ref.Path, "/"+actionEditPagePathKey+"/", "", 1)

	case ActionServeResource:
		return strings.Replace(r.URL.Path, "/"+actionEditPagePathKey+"/", "", 1)

	case ActionSavePage:
		return strings.Replace(r.URL.Path, "/"+actionSavePagePathKey+"/", "", 1)

	case ActionLoadPage:
		return strings.Replace(r.URL.Path, "/"+actionLoadPagePathKey+"/", "", 1)

	case ActionUploadImage:
		return strings.Replace(r.URL.Path, "/"+actionUploadImagePathKey+"/", "", 1)
	}

	return r.URL.Path
}

// getURLForAction retrieves the URL for a specific action
func getURLForAction(r *http.Request, a Action) url.URL {

	switch a {

	case ActionJsRequest:
		return url.URL{
			Host:   r.Host,
			Scheme: r.URL.Scheme,
			Path:   actionJsRequestPathKey}

	case ActionSavePage:
		return url.URL{
			Host:   r.Host,
			Scheme: r.URL.Scheme,
			Path:   actionSavePagePathKey}

	case ActionLoadPage:
		return url.URL{
			Host:   r.Host,
			Scheme: r.URL.Scheme,
			Path:   actionLoadPagePathKey}

	case ActionUploadImage:
		return url.URL{
			Host:   r.Host,
			Scheme: r.URL.Scheme,
			Path:   actionUploadImagePathKey}
	}

	return url.URL{
		Host:   r.Host,
		Scheme: r.URL.Scheme,
		Path:   "/"}
}

// getAction extracts request action
func getAction(r *http.Request) (a Action) {

	path := strings.Trim(r.URL.Path, "/")
	pathSplit := strings.Split(path, "/")

	if len(pathSplit) < 0 {
		return ActionGetPage
	}

	switch strings.ToLower(pathSplit[0]) {
	case actionJsRequestPathKey:
		return ActionJsRequest
	case actionSavePagePathKey:
		return ActionSavePage
	case actionLoadPagePathKey:
		return ActionLoadPage
	case actionUploadImagePathKey:
		return ActionUploadImage
	}

	return ActionServeResource
}
