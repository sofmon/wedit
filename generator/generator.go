// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const editorScriptTag = `<script async="" src="/!editor.js"></script>`

// Generator is used to generate static HTML websites through wedit web editor
type Generator struct {
	settings      Settings
	fileExplorer  FileExplorer
	publicHandler http.Handler
	Done          chan bool
}

// NewGenerator creates new wedit Generator
func NewGenerator(path string) (*Generator, error) {

	settings, err := NewSettings(path)
	if err != nil {
		return nil, err
	}

	fileExplorer := NewFileExplorer(settings)

	publicHandler := http.FileServer(http.Dir(settings.Folders.Public))

	done := make(chan bool)

	return &Generator{settings, fileExplorer, publicHandler, done}, nil
}

// Serve blocks the current goroutine, opens a port, listen and serves requests for wedit
func (g *Generator) Serve() error {

	// Editor settings
	generatorHandler := func(w http.ResponseWriter, r *http.Request) { handle(w, r, g) }
	http.HandleFunc("/!/", generatorHandler)
	http.HandleFunc("/!editor.js", generatorHandler)
	http.HandleFunc("/!save/", generatorHandler)
	http.HandleFunc("/!upload/", generatorHandler)
	endPoint := fmt.Sprintf("%v:%d", g.settings.Editor.Host, g.settings.Editor.Port)

	// Static website handling
	http.Handle("/", g.publicHandler)

	log.Printf("Start listening at http://%s/", endPoint)
	err := http.ListenAndServe(endPoint, nil)
	if err != nil {
		log.Printf("Unable to open editor HTTP server at '%v'. Error: %v\n", endPoint, err)
		return err
	}

	g.Done <- true

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
	// ErrMethodNotSupported is returned when the HTTP method is not supported
	ErrMethodNotSupported = errors.New("HTTP method not supported")

	// ErrEmptyPageSave is returned when save page request is empty
	ErrEmptyPageSave = errors.New("Page save request has empty page body")

	// ErrPageNotFound is returned when a page is not found
	ErrPageNotFound = errors.New("Page not found")
)

// Get requested page by subject
func (g *Generator) retrievePageForSubject(subject Subject) (*Page, error) {

	// Retrieve the page
	page, err := g.fileExplorer.ReadPageData(subject.Path())
	if err != nil && err != ErrPageNotFound {
		return nil, err
	}

	return &page, nil
}

func (g *Generator) handleServeResource(subject Subject) {

	path := subject.Path()

	// If there is an extension in the request, it will be served through the public handler as asset reqeust
	if strings.LastIndex(path, ".") > strings.LastIndex(path, "/") {
		subject.Request.URL.RawPath = strings.Replace(subject.Request.URL.RawPath, "/!/", "/", 1)
		subject.Request.URL.Path = strings.Replace(subject.Request.URL.Path, "/!/", "/", 1)
		g.publicHandler.ServeHTTP(subject.Response, subject.Request)
		return
	}

	if !strings.HasSuffix(subject.Request.URL.RawPath, "/") {
		http.Redirect(subject.Response, subject.Request, subject.Request.URL.RawPath+"/", http.StatusPermanentRedirect)
		return
	}

	// Otherwise, serve it as a new template with editor
	templateRow, err := g.fileExplorer.ReadPageTemplate(path)
	if err != nil {
		log.Printf("Unable to serve template on path '%v'. Error: %v", path, err)
		serveError(subject.Response, ErrPageNotFound)
		return
	}

	templateSplit := strings.Split(templateRow, "</body>")
	templateWithEditor := templateSplit[0] + editorScriptTag + "</body>"

	subject.Response.Write([]byte(templateWithEditor))

}

func (g *Generator) handleJsRequest(subject Subject) {

	// Retrieve page for the current request
	page, pageErr := g.retrievePageForSubject(subject)
	if pageErr != nil {
		log.Printf("Could not create page data. Path: %v; Error: %v;\n", subject.Path(), pageErr)
		serveError(subject.Response, pageErr)
		return
	}

	// Generate the wedit JS for the requested page
	preparedJs, err := getPreparedJs(subject, page)
	if err != nil {
		log.Printf("Could not prepare wedit.js. Key: %v; Error: %v;\n", subject.Path(), err)
		serveError(subject.Response, pageErr)
		return
	}

	// Returned the prepared site wedit JS
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

	err = g.fileExplorer.WritePageData(subject.Path(), *page)
	if err != nil {
		serveError(subject.Response, err)
	}

	// TODO: generate

}

func (g *Generator) handleUploadImage(subject Subject) {
	data, err := ioutil.ReadAll(subject.Request.Body)
	if err != nil {
		log.Printf("Could not save asset. Path: %v; Error: %v;\n", subject.Path(), err)
		serveError(subject.Response, err)
	}

	err = g.fileExplorer.WriteAsset(subject.Path(), data)
	if err != nil {
		log.Printf("Could not save asset. Path: %v; Error: %v;\n", subject.Path(), err)
		serveError(subject.Response, err)
	}
}

func handle(w http.ResponseWriter, r *http.Request, g *Generator) {

	w.Header().Add("ache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Pragma", "no-cache")
	w.Header().Add("Expires", "0")

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
