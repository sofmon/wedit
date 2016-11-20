// Copyright (c) 2016, Haralampi Staykov. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
package generator

import (
	"strings"

	"net/http"
	"net/url"
)

// Action to be performed by the current request
type Action int

const (
	// ActionServeResource is used when HTML or static resourse is requested
	ActionServeResource Action = iota
	// ActionJsRequest is used when the editor JS file is requested
	ActionJsRequest
	// ActionSavePage is used when page is beeing saved by the editor
	ActionSavePage
	// ActionLoadPage is used when page is beeing loaded by the editro
	ActionLoadPage
	// ActionUploadImage is used when image is uploaded by the editor
	ActionUploadImage
)

// Subject represents the wedit hcurrent request
type Subject struct {

	// Current HTTP request object
	Request *http.Request

	// HTTP response writer for the current request
	Response http.ResponseWriter

	// Host name of the varchar instance serving this page
	Host string

	// Unique path to the requested page
	Path string

	// Action to be taken on that page
	Action Action

	// A key targeted by the action
	Key string
}

const (
	actionJsRequestPathKey   = "editor.js"
	actionSavePagePathKey    = "save-page"
	actionLoadPagePathKey    = "load-page"
	actionUploadImagePathKey = "upload-image"
)

// GetURLForAction retrieves the URL for a specific action
func (subject Subject) GetURLForAction(action Action) url.URL {

	switch action {

	case ActionJsRequest:
		return url.URL{
			Host:   subject.Host,
			Scheme: subject.Request.URL.Scheme,
			Path:   actionJsRequestPathKey}

	case ActionSavePage:
		return url.URL{
			Host:   subject.Host,
			Scheme: subject.Request.URL.Scheme,
			Path:   "/!!/" + actionSavePagePathKey}

	case ActionLoadPage:
		return url.URL{
			Host:   subject.Host,
			Scheme: subject.Request.URL.Scheme,
			Path:   "/!!/" + actionLoadPagePathKey}

	case ActionUploadImage:
		return url.URL{
			Host:   subject.Host,
			Scheme: subject.Request.URL.Scheme,
			Path:   "/!!/" + actionUploadImagePathKey}
	}

	return url.URL{
		Host:   subject.Host,
		Scheme: subject.Request.URL.Scheme,
		Path:   "/"}
}

// NewSubject is used to create new Subject for a request
func NewSubject(request *http.Request, responseWriter http.ResponseWriter) Subject {

	var subject Subject

	subject.Request = request
	subject.Response = responseWriter

	subject.Host = request.Host

	path := strings.Trim(request.URL.Path, "/")
	pathSplit := strings.Split(path, "/")

	// The action to be taken by the server
	subject.Action = ActionServeResource
	if len(pathSplit) > 1 && pathSplit[0] == "!!" {
		// Action on the page
		switch strings.ToLower(pathSplit[1]) {
		case actionJsRequestPathKey:
			subject.Action = ActionJsRequest
			break
		case actionSavePagePathKey:
			subject.Action = ActionSavePage
			break
		case actionLoadPagePathKey:
			subject.Action = ActionLoadPage
			break
		case actionUploadImagePathKey:
			subject.Action = ActionUploadImage
			break
		}
	}

	// Get the key
	query := request.URL.Query()
	subject.Key = query.Get("k")

	// Get path to the current CMS-able page from the script referer
	subject.Path = query.Get("p")
	if subject.Path == "" {
		ref, err := url.Parse(request.Referer())
		if err == nil {
			subject.Path = ref.Path
		} else {
			subject.Path = "/"
		}
	}

	return subject
}
