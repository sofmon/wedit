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
	"strings"
)

const editorScriptTag = `<script async src="/!!/editor.js"></script>`

type Generator struct {
	settings      Settings
	fileExplorer  FileExplorer
	publicHandler http.Handler
}

func NewGenerator(path string) (*Generator, error) {

	settings, err := NewSettings(path)
	if err != nil {
		return nil, err
	}

	fileExplorer := NewFileExplorer(settings)

	publicHandler := http.FileServer(http.Dir(settings.Folders.Public))

	return &Generator{settings, fileExplorer, publicHandler}, nil
}

func (g *Generator) Serve() error {

	// Editor settings
	prefix := fmt.Sprintf("/%v/", g.settings.Editor.Prefix)
	generatorHandler := func(w http.ResponseWriter, r *http.Request) { handle(w, r, g) }
	http.HandleFunc(prefix, generatorHandler)
	http.HandleFunc("/!!/", generatorHandler)
	endPoint := fmt.Sprintf("%v:%d", g.settings.Editor.Host, g.settings.Editor.Port)

	// Static website handaling
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

	path := strings.Replace(subject.Request.URL.Path, "/"+g.settings.Editor.Prefix, "", 1)

	// If there is an extention in the request, it will be served through the public handler as asset reqeust
	if strings.LastIndex(path, ".") > strings.LastIndex(path, "/") {
		g.publicHandler.ServeHTTP(subject.Response, subject.Request)
		return
	}

	// Otherwise, serve it as a new template with editor
	tempalteRow, err := g.fileExplorer.ReadPageTemplate(path)
	if err != nil {
		log.Printf("Unable to serve template on path '%v'. Error: %v", path, err)
		serveError(subject.Response, ErrPageNotFound)
		return
	}

	templateSplit := strings.Split(tempalteRow, "</body>")
	templateWithEditor := templateSplit[0] + editorScriptTag + "</body>"

	subject.Response.Write([]byte(templateWithEditor))

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
