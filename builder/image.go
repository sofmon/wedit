package builder

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	// Support image codecs
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
	"golang.org/x/net/html"

	"github.com/sofmon/wedit/model"
)

func processImage(k model.Key, n *html.Node, p *model.Page) {

	img, found := p.Images[k]
	if !found {
		return
	}

	if !img.IsSet() {
		return
	}

	for i, a := range n.Attr {
		if strings.ToLower(a.Key) == "src" {
			n.Attr[i].Val = img.FileNameMin()
		}
		if strings.ToLower(a.Key) == "srcset" {
			n.Attr[i].Val = img.Srcset()
		}
	}
}

func processNodeForSrcset(page *model.Page, n *html.Node) error {

	if n.Data == "img" {

		var key model.Key
		var _, srcset string
		for _, a := range n.Attr {
			switch a.Key {
			case cfg.EditAttr:
				key = model.Key(a.Val)
			case "srcset":
				srcset = a.Val
			}
		}

		imgType := ""
		if img, found := page.Images[key]; found {
			imgType = img.Type
		}

		page.Images[key] = model.NewImage(key, imgType, srcset)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err := processNodeForSrcset(page, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateImagesSrcset(page *model.Page, path string) error {

	templateFile, err := findTemplateFile(path)
	if err != nil {
		return err
	}

	doc, err := openTemlateHTML(templateFile)
	if err != nil {
		return err
	}

	err = includesProcessNode(doc)
	if err != nil {
		return err
	}

	err = processNodeForSrcset(page, doc)
	if err != nil {
		return err
	}

	return nil
}

// WriteImage to the content and public folder
func WriteImage(pagePath string, key model.Key, name string, data []byte) (img model.Image, err error) {

	publicFolder := cfg.PublicFolder + pagePath
	os.MkdirAll(publicFolder, 0777)

	contentFolder := cfg.ContentFolder + pagePath
	os.MkdirAll(contentFolder, 0777)

	ext := path.Ext(name)
	imageType := strings.ToLower(strings.TrimLeft(ext, "."))

	file := contentFolder + string(key) + ext

	err = ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("✘ unable to save file %s; error: %v", file, err)
		return
	}

	page, err := ReadPageData(pagePath)
	if err != nil {
		log.Printf("✘ unable to read page data for path %s; error: %v", pagePath, err)
		return
	}

	found := false
	img, found = page.Images[key]
	if !found {
		log.Printf("✘ unable to find image with key '%s' in template for path %s", key, pagePath)
	}

	img.Type = imageType

	err = resizeImage(publicFolder, img, data)
	if err != nil {
		log.Printf("✘ unable to resize files; error: %v", err)
		return
	}

	return
}

func resizeImage(publicFolder string, img model.Image, data []byte) (err error) {

	var imgData image.Image
	switch img.Type {
	case "jpg", "jpeg":
		imgData, err = jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
	case "png":
		imgData, err = png.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
	case "gif":
		imgData, err = gif.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported image format %s", img.Type)
	}

	for _, w := range img.Width {

		rImgData := resize.Resize(uint(w), 0, imgData, resize.Lanczos3)

		fileName := publicFolder + img.FileNameW(w)

		out, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("unable to create file %s; error: %v", fileName, err)
		}

		// write new image to file
		switch img.Type {
		case "jpg", "jpeg":
			err = jpeg.Encode(out, rImgData, nil)
		case "png":
			err = png.Encode(out, rImgData)
		case "gif":
			err = gif.Encode(out, rImgData, nil)
		default:
			return fmt.Errorf("unsupported image format %s", img.Type)
		}

		out.Close()

		if err != nil {
			return fmt.Errorf("unable to resize file %s; error: %v", fileName, err)
		}
	}

	return nil
}
