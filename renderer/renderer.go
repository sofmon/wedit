// Package renderer generates HTML files in the `public` folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package renderer

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

// Renderer generated HTML files
type Renderer interface {
	RenderHTML(targetFilePath, templateFilePath string, page model.Page) error
}

type renderer struct {
	cfg Config
}

// NewRenderer creates new Renderer
func NewRenderer(cfg Config) Renderer {
	return &renderer{cfg}
}

// RenderHTML file based on template and page files
func (r *renderer) RenderHTML(targetFilePath, templateFilePath string, page model.Page) error {

	tplBytes, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		return fmt.Errorf("renderer: unable to read template file (%s) due to error: %v", templateFilePath, err)
	}

	doc, err := html.Parse(bytes.NewReader(tplBytes))
	if err != nil {
		return fmt.Errorf("renderer: unable to parse template file ('%s') as HTML due to error: %v", templateFilePath, err)
	}

	r.processNode(doc, &page)

	var targetBuff bytes.Buffer
	w := bufio.NewWriter(&targetBuff)

	err = html.Render(w, doc)
	if err != nil {
		return fmt.Errorf("renderer: unable to render HTML file (%s) due to error: %v", targetFilePath, err)
	}

	err = w.Flush()
	if err != nil {
		return fmt.Errorf("renderer: unable to render HTML file (%s) due to error: %v", targetFilePath, err)
	}

	err = ioutil.WriteFile(targetFilePath, targetBuff.Bytes(), 0777)
	if err != nil {
		return fmt.Errorf("renderer: unable to save HTML file (%s) due to error: %v", targetFilePath, err)
	}

	return nil
}

func (r *renderer) processNode(n *html.Node, page *model.Page) {
	for i, a := range n.Attr {
		if a.Key == r.cfg.EditAttr {
			k := model.Key(a.Val)
			if strings.ToLower(string(n.Data)) == "img" {
				r.processImage(k, n, page)
			} else {
				r.processElement(k, n, page)
			}
			if !r.cfg.KeepAttr {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}

		if a.Key == r.cfg.RepeatAttr {
			k := model.Key(a.Val)
			r.processRepeat(k, n, page)
			if !r.cfg.KeepAttr {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.processNode(c, page)
	}
}
