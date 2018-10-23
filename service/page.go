// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

const editorScriptTag = `<script async="" src="/!editor.js"></script>`

func (s *service) pageHandler(w http.ResponseWriter, r *http.Request) {

	path := getPathWithoutAction(r)
  builderConfig := s.bld.GetConfig()
	pageFile := "index"
	isPageFile := false
	if filepath.Ext(path) != "" {
		pageFile = strings.TrimSuffix(filepath.Base(path), filepath.Ext(filepath.Base(path)))
		isPageFile = true
		if !builderConfig.Contains(builderConfig.PageFiles, pageFile) {
			r.URL.RawPath = strings.Replace(r.URL.RawPath, "/!/", "/", 1)
			r.URL.Path = strings.Replace(r.URL.Path, "/!/", "/", 1)
			s.staticHandler.ServeHTTP(w, r)
			return
		}
	}

	if !strings.HasSuffix(r.URL.RawPath, "/") && !isPageFile {
		http.Redirect(w, r, r.URL.RawPath+"/", http.StatusPermanentRedirect)
		return
	}

	templateRaw, err := s.bld.ReadPageTemplate(path, pageFile)
	if err != nil {
		log.Printf("unable to serve template on path '%v'. Error: %v", path, err)
		http.NotFound(w, r)
		return
	}

	templateHTML, err := processEditLinks(templateRaw)
	if err != nil {
		log.Printf("unable to serve template on path '%v'. Error: %v", path, err)
		http.NotFound(w, r)
		return
	}

	templateSplit := strings.Split(templateHTML, "</body>")
	templateWithEditor := templateSplit[0] + editorScriptTag + "</body>"

	w.Write([]byte(templateWithEditor))
}

func processEditLinks(inHTML string) (outHTML string, err error) {

	doc, err := html.Parse(bytes.NewReader([]byte(inHTML)))

	processEditLinksNode(doc)

	var targetBuff bytes.Buffer
	w := bufio.NewWriter(&targetBuff)

	err = html.Render(w, doc)
	if err != nil {
		err = fmt.Errorf("unable to render HTML due to error: %v", err)
		return
	}

	err = w.Flush()
	if err != nil {
		err = fmt.Errorf("unable to render HTML due to error: %v", err)
		return
	}

	outHTML = string(targetBuff.Bytes())

	return
}

func processEditLinksNode(n *html.Node) {

	if n.Data == "a" {
		for i, a := range n.Attr {
			if a.Key == "href" &&
				strings.HasPrefix(a.Val, "/") {
				n.Attr[i].Val = "/!" + a.Val
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processEditLinksNode(c)
	}
}
