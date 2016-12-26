// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

import (
	"strings"

	"net/http"
	"net/url"
)

// Action to be performed by the current request
type Action int

const (
	// ActionServeResource is used when HTML or static resource is requested
	ActionServeResource Action = iota
	// ActionJsRequest is used when the editor JS file is requested
	ActionJsRequest
	// ActionSavePage is used when page is being saved by the editor
	ActionSavePage
	// ActionLoadPage is used when page is being loaded by the editor
	ActionLoadPage
	// ActionUploadImage is used when image is uploaded by the editor
	ActionUploadImage
)

// Subject represents the wedit current request
type Subject struct {

	// Current HTTP request object
	Request *http.Request

	// HTTP response writer for the current request
	Response http.ResponseWriter

	// Host name of the wedit instance serving this page
	Host string

	// Action to be taken on that page
	Action Action
}

const (
	actionJsRequestPathKey   = "!editor.js"
	actionEditPagePathKey    = "!"
	actionSavePagePathKey    = "!save"
	actionLoadPagePathKey    = "!load"
	actionUploadImagePathKey = "!upload"
)

// Path retrieves the path of the requested resource excluding the wedit parameters
func (subject Subject) Path() string {

	switch subject.Action {

	case ActionJsRequest:
		ref, err := url.Parse(subject.Request.Referer())
		if err != nil {
			return "/" // TODO: evaluate
		}
		return strings.Replace(ref.Path, "/"+actionEditPagePathKey+"/", "", 1)

	case ActionServeResource:
		return strings.Replace(subject.Request.URL.Path, "/"+actionEditPagePathKey+"/", "", 1)

	case ActionSavePage:
		return strings.Replace(subject.Request.URL.Path, "/"+actionSavePagePathKey+"/", "", 1)

	case ActionLoadPage:
		return strings.Replace(subject.Request.URL.Path, "/"+actionLoadPagePathKey+"/", "", 1)

	case ActionUploadImage:
		return strings.Replace(subject.Request.URL.Path, "/"+actionUploadImagePathKey+"/", "", 1)
	}

	return subject.Request.URL.Path
}

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
			Path:   actionSavePagePathKey}

	case ActionLoadPage:
		return url.URL{
			Host:   subject.Host,
			Scheme: subject.Request.URL.Scheme,
			Path:   actionLoadPagePathKey}

	case ActionUploadImage:
		return url.URL{
			Host:   subject.Host,
			Scheme: subject.Request.URL.Scheme,
			Path:   actionUploadImagePathKey}
	}

	return url.URL{
		Host:   subject.Host,
		Scheme: subject.Request.URL.Scheme,
		Path:   "/"}
}

// NewSubject creates new Subject for a request
func NewSubject(request *http.Request, responseWriter http.ResponseWriter) Subject {

	var subject Subject

	subject.Request = request
	subject.Response = responseWriter

	subject.Host = request.Host

	path := strings.Trim(request.URL.Path, "/")
	pathSplit := strings.Split(path, "/")

	// The action to be taken by the server
	subject.Action = ActionServeResource

	if len(pathSplit) > 0 {

		// Action on the page
		switch strings.ToLower(pathSplit[0]) {
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

	return subject
}
