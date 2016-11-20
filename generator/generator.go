// Package generator by www.sofmon.com
// Copyright (c) 2016, Haralampi Staykov. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
package generator

import (
	"errors"

	"fmt"
	"io"
	"log"
	"net/http"
)

type Generator struct {
	settings      Settings
	fileExplorer  FileExplorer
	publicHandler http.Handler
	dataHandler   http.Handler
}

func NewGenerator(path string) (*Generator, error) {

	settings, err := NewSettings(path)
	if err != nil {
		return nil, err
	}

	fileExplorer := NewFileExplorer(settings)

	return &Generator{settings, fileExplorer, nil, nil}, nil
}

func (g *Generator) Serve() error {

	// Editor settings
	prefix := fmt.Sprintf("/%v/", g.settings.Editor.Prefix)
	http.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) { handle(w, r, g) })
	endPoint := fmt.Sprintf("%v:%d", g.settings.Editor.Host, g.settings.Editor.Port)
	g.dataHandler = http.FileServer(http.Dir(g.settings.Folders.Data)) // to be used when serving static assets

	// Static website
	g.publicHandler = http.FileServer(http.Dir(g.settings.Folders.Public))
	http.Handle("/", g.publicHandler)

	log.Printf("Start listening at http://%s/", endPoint)
	err := http.ListenAndServe(endPoint, nil)
	if err != nil {
		log.Printf("Unable to open editor HTTP server at '%v'. Error: %v\n", endPoint, err)
		return err
	}

	return nil
}

func serveError(w http.ResponseWriter, err error) {
	switch err {
	case ErrPageNotFound:
		w.WriteHeader(http.StatusNotFound)
		break
	case ErrMethodNotSupported:
		w.WriteHeader(http.StatusMethodNotAllowed)
		break
	case ErrEmptyPageSave:
		w.WriteHeader(http.StatusBadRequest)
		break
	default:
		w.WriteHeader(http.StatusInternalServerError)
		break
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, err.Error())
	io.WriteString(w, "Internal Server Error")
}

var (
	ErrMethodNotSupported = errors.New("HTTP method not supported")
	ErrEmptyPageSave      = errors.New("Page save request has empty page body")
	ErrPageNotFound       = errors.New("Page not found")
)

// Get requested page by subject
func (g *Generator) retrievePageForSubject(subject Subject) (*Page, error) {

	// Retrieve the page
	page, err := NewEmptyPage(), error(nil) //TODO: LoadPageFromDatastore(subject)
	if err != nil && err != ErrPageNotFound {
		return nil, err
	}

	if page != nil {
		return page, nil
	}

	// No page found, creation new one
	log.Printf("First request for page: Path: %v;\n", subject.Path)

	return NewEmptyPage(), nil
}

func (g *Generator) handleServeResource(subject Subject) {

	log.Println(subject.Request.URL.Path)

	// If it is not an html file serve as static resource from data folder
	g.dataHandler.ServeHTTP(subject.Response, subject.Request)
}

func (g *Generator) handleJsRequest(subject Subject) {

	// Retrieve page for the current request
	page, pageErr := g.retrievePageForSubject(subject)
	if pageErr != nil {
		log.Printf("Could not create page data. Path: %v; Error: %v;\n", subject.Path, pageErr)
		serveError(subject.Response, pageErr)
		return
	}

	// Generate the varchar JS for the requested page
	preparedJs, err := getPreparedJs(subject, page)
	if err != nil {
		log.Printf("Could not prepare varchar.js. Key: %v; Error: %v;\n", subject.Path, err)
		serveError(subject.Response, pageErr)
		return
	}

	// Returned the prepared site varchar JS
	subject.Response.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	subject.Response.Write([]byte(*preparedJs))

}

func (g *Generator) handleSavePage(subject Subject) {

	page, err := ReadPageFromJson(subject.Request.Body)
	if err != nil {
		serveError(subject.Response, err)
		return
	}

	if page == nil {
		serveError(subject.Response, ErrEmptyPageSave)
		return
	}

	err = error(nil) //TODO: SavePageToDatastore(subject, page)
	if err != nil {
		serveError(subject.Response, err)
	}

	// TODO: generate

}

func (g *Generator) handleUploadImage(subject Subject) {
	// HERE
}

func handle(w http.ResponseWriter, r *http.Request, g *Generator) {

	subject := NewSubject(r, w)

	switch subject.Action {

	case ActionServeResource:
		g.handleServeResource(subject)
		return

	case ActionJsRequest:
		g.handleJsRequest(subject)
		return

	case ActionSavePage:
		g.handleSavePage(subject)
		return

	case ActionUploadImage:
		g.handleUploadImage(subject)
		return
	}

	serveError(w, ErrMethodNotSupported)
}
