package renderer

import (
	"bufio"
	"bytes"
	"strings"

	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

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
