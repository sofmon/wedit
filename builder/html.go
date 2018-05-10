package builder

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

func openTemlateHTML(templateFilePath string) (doc *html.Node, err error) {

	tplBytes, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		err = fmt.Errorf("unable to read template file (%s) due to error: %v", templateFilePath, err)
		return
	}

	doc, err = html.Parse(bytes.NewReader(tplBytes))
	if err != nil {
		err = fmt.Errorf("unable to parse template file ('%s') as HTML due to error: %v", templateFilePath, err)
		return
	}

	return
}

func openIncludeHTML(templateFilePath string, parent *html.Node) (nodes []*html.Node, err error) {

	tplBytes, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		err = fmt.Errorf("unable to read include file (%s) due to error: %v", templateFilePath, err)
		return
	}

	nodes, err = html.ParseFragment(bytes.NewReader(tplBytes), parent)
	if err != nil {
		err = fmt.Errorf("unable to parse include file ('%s') as HTML due to error: %v", templateFilePath, err)
		return
	}

	return
}

func saveTargetHTML(targetFilePath string, doc *html.Node) (err error) {

	var targetBuff bytes.Buffer
	w := bufio.NewWriter(&targetBuff)

	err = html.Render(w, doc)
	if err != nil {
		return fmt.Errorf("unable to render HTML file (%s) due to error: %v", targetFilePath, err)
	}

	err = w.Flush()
	if err != nil {
		return fmt.Errorf("unable to render HTML file (%s) due to error: %v", targetFilePath, err)
	}

	err = ioutil.WriteFile(targetFilePath, targetBuff.Bytes(), 0777)
	if err != nil {
		return fmt.Errorf("unable to save HTML file (%s) due to error: %v", targetFilePath, err)
	}

	return
}

func (b *builder) renderHTML(targetFilePath, templateFilePath string, page model.Page) error {

	doc, err := openTemlateHTML(targetFilePath)
	if err != nil {
		return err
	}

	b.includesProcessNode(doc)
	b.renderProcessNode(doc, &page)

	err = saveTargetHTML(targetFilePath, doc)
	if err != nil {
		return err
	}

	return nil
}

func (b *builder) renderProcessNode(n *html.Node, page *model.Page) {
	for i, a := range n.Attr {
		if a.Key == b.cfg.EditAttr {
			k := model.Key(a.Val)
			if strings.ToLower(string(n.Data)) == "img" {
				processImage(k, n, page)
			} else {
				processElement(k, n, page)
			}
			if !b.cfg.KeepWeditAttrs {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}

		if a.Key == b.cfg.RepeatAttr {
			k := model.Key(a.Val)
			b.processRepeat(k, n, page)
			if !b.cfg.KeepWeditAttrs {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.renderProcessNode(c, page)
	}
}

func (b *builder) prepareIncludes(targetFilePath string) (outHTML []byte, err error) {

	doc, err := openTemlateHTML(targetFilePath)
	if err != nil {
		return
	}

	err = b.includesProcessNode(doc)
	if err != nil {
		return
	}

	var targetBuff bytes.Buffer
	w := bufio.NewWriter(&targetBuff)

	err = html.Render(w, doc)
	if err != nil {
		err = fmt.Errorf("renderer: unable to render HTML due to error: %v", err)
		return
	}

	err = w.Flush()
	if err != nil {
		err = fmt.Errorf("renderer: unable to render HTML due to error: %v", err)
		return
	}

	outHTML = targetBuff.Bytes()

	return
}

func (b *builder) includesProcessNode(n *html.Node) error {
	for i, a := range n.Attr {
		if a.Key == b.cfg.IncludeAttr {

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				n.RemoveChild(c)
			}

			if !b.cfg.KeepWeditAttrs {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}

			ics, err := b.getInclude(a.Val, n)
			if err != nil {
				return fmt.Errorf("unable to include '%v' doe to an error: %v", a.Val, err)
			}

			for _, ic := range ics {
				n.AppendChild(ic)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := b.includesProcessNode(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *builder) getInclude(path string, parent *html.Node) (nodes []*html.Node, err error) {

	filePath := strings.TrimSuffix(b.cfg.TemplateFolder, "/") + "/" + strings.TrimPrefix(path, "/")

	return openIncludeHTML(filePath, parent)
}
