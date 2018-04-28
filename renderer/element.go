package renderer

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/russross/blackfriday"

	"github.com/sofmon/wedit/model"
)

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
