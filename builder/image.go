package builder

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

func processImage(k model.Key, n *html.Node, p *model.Page) {

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

func (b *builder) WriteImage(path string, key, name string, data []byte) (err error) {

	publicFolder := b.cfg.PublicFolder + path
	os.MkdirAll(publicFolder, 0777)

	contentFolder := b.cfg.ContentFolder + path
	os.MkdirAll(contentFolder, 0777)

	file := contentFolder + name

	err = ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("unable to save file %s; error: %v", file, err)
		return err
	}

	// process public folder

	file = publicFolder + name

	err = ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("unable to save file %s; error: %v", file, err)
		return err
	}

	return
}
