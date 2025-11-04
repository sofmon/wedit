package builder

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html"

	"wedit/model"
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

func renderHTML(targetFilePath, templateFilePath string, page model.Page) error {

	doc, err := openTemlateHTML(templateFilePath)
	if err != nil {
		return err
	}

	err = includesProcessNode(doc)
	if err != nil {
		return err
	}

	renderProcessNode(doc, &page)

	err = saveTargetHTML(targetFilePath, doc)
	if err != nil {
		return err
	}

	return nil
}

func renderProcessNode(n *html.Node, page *model.Page) {
	for i, a := range n.Attr {
		if a.Key == cfg.EditAttr {
			k := model.Key(a.Val)
			if strings.ToLower(string(n.Data)) == "img" {
				processImage(k, n, page)
			} else {
				processElement(k, n, page)
			}
			if !cfg.KeepWeditAttrs {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}

		if a.Key == cfg.RepeatAttr {
			k := model.Key(a.Val)
			processRepeat(k, n, page)
			if !cfg.KeepWeditAttrs {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}

		if a.Key == cfg.ClassAttr {
			k := model.Key(a.Val)
			processClass(k, n, page)
			if !cfg.KeepWeditAttrs {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderProcessNode(c, page)
	}
}

func prepareIncludes(path string) (outHTML []byte, err error) {

	templateFilePath, err := findTemplateFile(path)
	if err != nil {
		return
	}

	doc, err := openTemlateHTML(templateFilePath)
	if err != nil {
		return
	}

	err = includesProcessNode(doc)
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

func includesProcessNode(n *html.Node) error {
	for i, a := range n.Attr {
		if a.Key == cfg.IncludeAttr {

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				n.RemoveChild(c)
			}

			if !cfg.KeepWeditAttrs {
				n.Attr = append(n.Attr[:i], n.Attr[i+1:]...)
			}

			ics, err := getInclude(a.Val, n)
			if err != nil {
				return fmt.Errorf("unable to include '%v' doe to an error: %v", a.Val, err)
			}

			for _, ic := range ics {
				n.AppendChild(ic)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := includesProcessNode(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func getInclude(path string, parent *html.Node) (nodes []*html.Node, err error) {

	filePath := cfg.TemplateFolder + path

	return openIncludeHTML(filePath, parent)
}

func findNodeByKey(n *html.Node, key model.Key) *html.Node {

	for _, a := range n.Attr {
		if a.Key == cfg.EditAttr {
			k := model.Key(a.Val)
			if k == key {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		cn := findNodeByKey(c, key)
		if cn != nil {
			return cn
		}
	}

	return nil
}
