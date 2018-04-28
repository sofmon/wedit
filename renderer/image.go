package renderer

import (
	"strings"

	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

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
