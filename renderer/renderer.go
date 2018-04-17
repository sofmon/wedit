// Package renderer generates HTML files in the `public` folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package renderer

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

// Renderer generated HTML files
type Renderer interface {
	RenderHTML(targetFilePath, dataFilePath, templateFilePath string) error
}

type renderer struct {
	cfg Config
}

// NewRenderer creates new Renderer
func NewRenderer(cfg Config) Renderer {
	return &renderer{cfg}
}

// RenderHTML file based on template and page files
func (r *renderer) RenderHTML(targetFilePath, dataFilePath, templateFilePath string) error {

	dataBytes, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		return fmt.Errorf("renderer: unable to read data file (%s) due to error: %v", dataFilePath, err)
	}

	var page model.Page
	err = json.Unmarshal(dataBytes, &page)
	if err != nil {
		return fmt.Errorf("renderer: unable to parse data file (%s) as JSON due to error: %v", dataFilePath, err)
	}

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

func (r *renderer) processElement(k model.Key, n *html.Node, p *model.Page) {

	element, found := p.Elements[k]
	if !found {
		return
	}

	text := string(blackfriday.Run([]byte(element.Text)))
	if strings.Index(text, "<p>") == strings.LastIndex(text, "<p>") {
		text = strings.Replace(text, "<p>", "", 1)
		text = strings.Replace(text, "</p>", "", 1)
	}

	parsedMd, err := html.ParseFragment(strings.NewReader(text), n)
	if err != nil {
		return
	}

	for c := n.FirstChild; c != nil; c = n.FirstChild {
		n.RemoveChild(c)
	}

	for _, c := range parsedMd {
		n.AppendChild(c)
	}
}

func (r *renderer) processImage(k model.Key, n *html.Node, p *model.Page) {

	image, found := p.Images[k]
	if !found {
		return
	}

	for i, a := range n.Attr {
		if strings.ToLower(a.Key) != "src" {
			continue
		}
		n.Attr[i].Val = image.Src
		break
	}
}

func (r *renderer) processRepeat(k model.Key, n *html.Node, p *model.Page) {

	repeat, found := p.Repeats[k]
	if !found {
		return
	}

	before := true
	for _, k := range repeat.CopyKeys {

		if k == repeat.Key {
			before = false
			continue
		}

		cn := r.cloneNode(n)

		for i, a := range cn.Attr {
			if strings.ToLower(a.Key) == r.cfg.RepeatAttr {
				cn.Attr = append(cn.Attr[:i], cn.Attr[i+1:]...)
				break
			}
		}

		r.updateVarValues(cn, k)
		r.processNode(cn, p)

		if before {
			n.Parent.InsertBefore(cn, n)
		} else {
			if n.NextSibling != nil {
				n.Parent.InsertBefore(cn, n.NextSibling)
			} else {
				n.Parent.AppendChild(cn)
			}
		}
	}
}

func (r *renderer) cloneNode(n *html.Node) *html.Node {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err := html.Render(w, n)
	if err != nil {
		return nil
	}

	err = w.Flush()
	if err != nil {
		return nil
	}

	cn, err := html.ParseFragment(bytes.NewReader(b.Bytes()), n.Parent)
	if err != nil {
		return nil
	}

	if len(cn) < 1 {
		return nil
	}

	return cn[0]
}

func (r *renderer) updateVarValues(n *html.Node, k model.Key) {

	for i, a := range n.Attr {
		if strings.ToLower(a.Key) == r.cfg.EditAttr {
			n.Attr[i].Val = a.Val + string(k)
			break
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		r.updateVarValues(c, k)
	}
}
