package builder

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/russross/blackfriday"

	"wedit/model"
)

func processElement(k model.Key, n *html.Node, p *model.Page) {

	element, found := p.Elements[k]
	if !found {
		return
	}

	text := string(blackfriday.MarkdownCommon([]byte(element.Text)))

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
