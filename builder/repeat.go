package builder

import (
	"bufio"
	"bytes"
	"strings"

	"golang.org/x/net/html"

	"wedit/model"
)

func processRepeat(k model.Key, n *html.Node, p *model.Page) {

	repeat, found := p.Repeats[k]
	if !found {
		return
	}

	before := true
	latestNode := n
	for _, k := range repeat.CopyKeys {

		if k == repeat.Key {
			before = false
			continue
		}

		cn := cloneNode(n)

		for i, a := range cn.Attr {
			if strings.ToLower(a.Key) == cfg.RepeatAttr {
				cn.Attr = append(cn.Attr[:i], cn.Attr[i+1:]...)
				break
			}
		}

		updateVarValues(cn, k)
		renderProcessNode(cn, p)

		if before {
			n.Parent.InsertBefore(cn, n)
		} else {
			if latestNode.NextSibling != nil {
				n.Parent.InsertBefore(cn, latestNode.NextSibling)
			} else {
				n.Parent.AppendChild(cn)
			}
			latestNode = cn
		}
	}
}

func cloneNode(n *html.Node) *html.Node {

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

func updateVarValues(n *html.Node, k model.Key) {

	for i, a := range n.Attr {
		if strings.ToLower(a.Key) == cfg.EditAttr {
			n.Attr[i].Val = a.Val + string(k)
			break
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		updateVarValues(c, k)
	}
}
