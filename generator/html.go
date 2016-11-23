package generator

import (
	"bufio"

	"bytes"
	"io/ioutil"

	"strings"

	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
)

func writePageAsHtml(indexFile, templateFile string, page *Page) error {

	data, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		processNode(n, page)
	}
	f(doc)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err = html.Render(w, doc)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(indexFile, b.Bytes(), 0777)
	if err != nil {
		return err
	}

	return nil
}

func processNode(n *html.Node, page *Page) {
	for i, a := range n.Attr {
		if a.Key == "data-var" {
			if strings.ToLower(string(n.Data)) == "img" {
				processImage(a.Val, n, page)
			} else {
				processElement(a.Val, n, page)
			}
			n.Attr = append(n.Attr[:i], n.Attr[i+1:]...) // Remove data attributes for final render
		}

		if a.Key == "data-var-repeat" {
			processRepeat(a.Val, n, page)
			n.Attr = append(n.Attr[:i], n.Attr[i+1:]...) // Remove data attributes for final render
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processNode(c, page)
	}
}

func processElement(k string, n *html.Node, p *Page) {
	var element *Element
	for _, e := range p.Elements {
		if e.Key == k {
			element = &e
			break
		}
	}

	if element == nil {
		return
	}

	text := string(blackfriday.MarkdownBasic([]byte(element.Text)))
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

func processImage(k string, n *html.Node, p *Page) {
	var image *Image
	for _, i := range p.Images {
		if i.Key == k {
			image = &i
			break
		}
	}

	if image == nil {
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

func processRepeat(k string, n *html.Node, p *Page) {
	var repeat *Repeat
	for _, r := range p.Repeats {
		if r.Key == k {
			repeat = &r
			break
		}
	}

	if repeat == nil {
		return
	}

	before := true
	for _, k := range strings.Split(repeat.CopyKeys, ",") {

		if k == repeat.Key {
			before = false
			continue
		}

		cn := cloneNode(n)

		for i, a := range cn.Attr {
			if strings.ToLower(a.Key) == "data-var-repeat" {
				cn.Attr = append(cn.Attr[:i], cn.Attr[i+1:]...)
				break
			}
		}

		updateVarValues(cn, k)
		processNode(cn, p)

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

func cloneNode(n *html.Node) *html.Node {
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

func updateVarValues(n *html.Node, k string) {

	for i, a := range n.Attr {
		if strings.ToLower(a.Key) == "data-var" {
			n.Attr[i].Val = a.Val + k
			break
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		updateVarValues(c, k)
	}
}
