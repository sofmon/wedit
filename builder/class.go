package builder

import (
	"golang.org/x/net/html"

	"wedit/model"
)

func processClass(k model.Key, n *html.Node, p *model.Page) {

	class, found := p.Classes[k]
	if !found {
		return
	}

	// If the class value is empty, don't modify the class attribute
	if class.Value == "" {
		return
	}

	// Find and update the class attribute
	for i := range n.Attr {
		if n.Attr[i].Key == "class" {
			// Replace the entire class attribute with the selected value
			n.Attr[i].Val = class.Value
			return
		}
	}

	// If no class attribute exists, add one
	n.Attr = append(n.Attr, html.Attribute{
		Key: "class",
		Val: class.Value,
	})
}
