package builder

import (
	"bufio"
	"bytes"
	"strings"

	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

func (b *builder) processRepeat(k model.Key, n *html.Node, p *model.Page) {

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

		cn := b.cloneNode(n)

		for i, a := range cn.Attr {
			if strings.ToLower(a.Key) == b.cfg.RepeatAttr {
				cn.Attr = append(cn.Attr[:i], cn.Attr[i+1:]...)
				break
			}
		}

		b.updateVarValues(cn, k)
		b.renderProcessNode(cn, p)

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

func (b *builder) cloneNode(n *html.Node) *html.Node {

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	err := html.Render(w, n)
	if err != nil {
		return nil
	}

	err = w.Flush()
	if err != nil {
		return nil
	}

	cn, err := html.ParseFragment(bytes.NewReader(buf.Bytes()), n.Parent)
	if err != nil {
		return nil
	}

	if len(cn) < 1 {
		return nil
	}

	return cn[0]
}

func (b *builder) updateVarValues(n *html.Node, k model.Key) {

	for i, a := range n.Attr {
		if strings.ToLower(a.Key) == b.cfg.EditAttr {
			n.Attr[i].Val = a.Val + string(k)
			break
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		b.updateVarValues(c, k)
	}
}
