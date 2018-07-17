// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

// Image represents an image on a page
type Image struct {
	Key          Key       `json:"k"`
	Type         string    `json:"t"`
	Width        []int     `json:"w"`
	PixelDencity []float64 `json:"x"`
}

// FileName for the image source without descriptors
func (img Image) FileName() string {
	return fmt.Sprintf("%s.%s", img.Key, img.Type)
}

// FileNameW for the image source with width descriptors
func (img Image) FileNameW(w int) string {
	return fmt.Sprintf("%s-%dw.%s", img.Key, w, img.Type)
}

// FileNameX for the image source with pixel density descriptors
func (img Image) FileNameX(pd float64) string {
	return fmt.Sprintf("%s-%.2fx%s", img.Key, pd, img.Type)
}

// FileNameMin for the image source as the smallest image
func (img Image) FileNameMin() string {

	if len(img.Width) > 0 {
		m := img.Width[0]
		for _, e := range img.Width {
			if e < m {
				m = e
			}
		}
		return img.FileNameW(m)
	}

	if len(img.PixelDencity) > 0 {
		m := img.PixelDencity[0]
		for _, e := range img.PixelDencity {
			if e < m {
				m = e
			}
		}
		return img.FileNameX(m)
	}

	return img.FileName()
}

// Srcset for the image srcset attribute
func (img Image) Srcset() string {

	buf := bytes.Buffer{}

	for _, w := range img.Width {
		buf.WriteString(fmt.Sprintf("%s %dw", img.FileNameW(w), w))
	}

	for _, x := range img.PixelDencity {
		buf.WriteString(fmt.Sprintf("%s %.2fx", img.FileNameX(x), x))
	}

	return buf.String()
}

/*
http://w3c.github.io/html/semantics-embedded-content.html#element-attrdef-img-srcset

The src attribute must be present, and must contain a valid non-empty URL potentially surrounded by spaces referencing a non-interactive, optionally animated, image resource that is neither paged nor scripted.

The srcset attribute may also be present. If present, its value must consist of one or more image candidate strings, each separated from the next by a U+002C COMMA character (,). If an image candidate string contains no descriptors and no space characters after the URL, the following image candidate string, if there is one, must begin with one or more space characters.

An image candidate string consists of the following components, in order, with the further restrictions described below this list:
1. Zero or more space characters.
2. A valid non-empty URL that does not start or end with a U+002C COMMA character (,), referencing a non-interactive, optionally animated, image resource that is neither paged nor scripted.
3. Zero or more space characters.
4. Zero or one of the following:
  - A width descriptor, consisting of: a space character, a valid non-negative integer giving a number greater than zero representing the width descriptor value, and a U+0077 LATIN SMALL LETTER W character.
  - A pixel density descriptor, consisting of: a space character, a valid floating-point number giving a number greater than zero representing the pixel density descriptor value, and a U+0078 LATIN SMALL LETTER X character.
5. Zero or more space characters.
*/

var imageCandidatesExtract = regexp.MustCompile(`([\w\d\.]+)[ ]+([\d\w\.]+[wx])`)

// NewImage from key and srcset
// Uploaded images has name derivered from the key
func NewImage(key Key, fileType string, srcset string) (img Image) {

	img.Key = key
	img.Type = fileType

	if srcset == "" {
		return
	}

	data := imageCandidatesExtract.FindAllStringSubmatch(srcset, -1)
	for _, kv := range data {
		_, desc := kv[1], kv[2]

		switch desc[len(desc)-1] {
		case 'w':
			if wd, err := strconv.Atoi(desc[:len(desc)-1]); err == nil {
				img.Width = append(img.Width, wd)
			}
		case 'x':
			if pd, err := strconv.ParseFloat(desc[:len(desc)-1], 64); err == nil {
				img.PixelDencity = append(img.PixelDencity, pd)
			}
		}
	}

	return
}
