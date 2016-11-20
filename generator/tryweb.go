package generator

/*
import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func WriteConvertedWebPage(w http.ResponseWriter, pageUrl string) error {

	targetUrl, err := url.Parse(pageUrl)
	if err != nil {
		return err
	}

	client := urlfetch.Client()
	resp, err := client.Get(pageUrl)
	if err != nil {
		return err
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	i := 0

	var f func(*html.Node)
	f = func(n *html.Node) {

		for i, a := range n.Attr {
			if a.Key == "src" || a.Key == "href" || a.Key == "srcset" {
				u, err := url.Parse(a.Val)
				if err == nil && !u.IsAbs() && u.Host == "" && u.Path != "/__try__" {
					u.Host = targetUrl.Host
					u.Scheme = targetUrl.Scheme
					//context.Infof("Converting %v -> %v;", a.Val, u.String())
					n.Attr[i].Val = u.String()
				}
			}
		}

		if n.Type == html.ElementNode && string(n.Data) == "body" {
			varcharScript, err := html.ParseFragment(strings.NewReader("<script async src=\"/__try__\"></script>"), n)
			if err != nil {
				return
			}
			varcharScript = varcharScript

			for _, c := range varcharScript {
				n.AppendChild(c)
			}
		}

		if n.Type == html.ElementNode &&
			string(n.Data) != "script" &&
			string(n.Data) != "style" &&
			n.FirstChild != nil &&
			n.FirstChild.Type == html.TextNode &&
			n.FirstChild.NextSibling == nil &&
			strings.Trim(string(n.FirstChild.Data), " ") != "" {

			n.Attr = append(n.Parent.Attr, html.Attribute{Key: "data-var", Val: fmt.Sprintf("v-%d", i)})
			i++
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	html.Render(w, doc)

	return nil
}
*/
