// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

const (
	keyPageHost    = "h"
	keyPageSite    = "s"
	keyPagePath    = "p"
	keyPageTitle   = "t"
	keyPageElement = "e"
	keyPageRepeat  = "r"
	keyPageImages  = "i"

	keyElementKey  = "k"
	keyElementText = "t"

	keyRepeatKey      = "k"
	keyRepeatCopyKeys = "c"

	keyImageKey = "k"
	keyImageSrc = "s"
)

func WritePageAsEscapedJson(path, host string, page *Page, buffer *bytes.Buffer) error {

	err := writeBuffer(buffer,
		"{\\\"", keyPageSite, "\\\":\\\"", "",
		"\\\",\\\"", keyPagePath, "\\\":\\\"", path,
		"\\\",\\\"", keyPageHost, "\\\":\\\"", host,
		"\\\",\\\"", keyPageTitle, "\\\":\\\"", page.Title,
		"\\\",\\\"", keyPageElement, "\\\":[")
	if err != nil {
		return err
	}

	for index, element := range page.Elements {

		if index > 0 {
			buffer.WriteRune(',')
		}

		err := WriteElementAsEscapedJson(&element, buffer)
		if err != nil {
			return err
		}
	}

	err = writeBuffer(buffer, "],\\\"", keyPageRepeat, "\\\":[")
	if err != nil {
		return err
	}

	for index, repeat := range page.Repeats {

		if index > 0 {
			buffer.WriteRune(',')
		}

		err := WriteRepeatAsEscapedJson(&repeat, buffer)
		if err != nil {
			return err
		}
	}

	err = writeBuffer(buffer, "],\\\"", keyPageImages, "\\\":[")
	if err != nil {
		return err
	}

	for index, image := range page.Images {

		if index > 0 {
			buffer.WriteRune(',')
		}

		err := WriteImageAsEscapedJson(&image, buffer)
		if err != nil {
			return err
		}
	}

	_, err = buffer.WriteString("]}")
	if err != nil {
		return err
	}

	return nil
}

func WriteElementAsEscapedJson(element *Element, buffer *bytes.Buffer) error {

	escaped := strings.Replace(element.Text, "'", "\\'", -1)
	escaped = strings.Replace(escaped, "\n", "<br>", -1)
	escaped = strings.Replace(escaped, "\"", "<q>", -1)

	err := writeBuffer(buffer,
		"{\\\"", keyElementKey, "\\\":\\\"", element.Key,
		"\\\",\\\"", keyElementText, "\\\":\\\"", escaped,
		"\\\"}")
	if err != nil {
		return err
	}
	return nil
}

func WriteRepeatAsEscapedJson(repeat *Repeat, buffer *bytes.Buffer) error {
	err := writeBuffer(buffer,
		"{\\\"", keyRepeatKey, "\\\":\\\"", repeat.Key,
		"\\\",\\\"", keyRepeatCopyKeys, "\\\":\\\"", repeat.CopyKeys,
		"\\\"}")
	if err != nil {
		return err
	}

	return nil
}

func WriteImageAsEscapedJson(image *Image, buffer *bytes.Buffer) error {
	err := writeBuffer(buffer,
		"{\\\"", keyImageKey, "\\\":\\\"", image.Key,
		"\\\",\\\"", keyImageSrc, "\\\":\\\"", image.Src,
		"\\\"}")
	if err != nil {
		return err
	}
	return nil
}

type PageMessage struct {
	T string           // Title
	E []ElementMessage // Elements
	R []RepeatMessage  // Repeats
	I []ImageMessage   // Repeats
}

type ElementMessage struct {
	K string // Key
	T string // Text
}

type RepeatMessage struct {
	K string // Key
	C string // CopyKeys
}

type ImageMessage struct {
	K string // Key
	S string // Src
}

func pageToMessage(p Page) (pm PageMessage) {

	elementToMessage := func(e Element) (em ElementMessage) {
		em.K = e.Key
		em.T = e.Text
		return
	}

	repeatToMessage := func(r Repeat) (rm RepeatMessage) {
		rm.K = r.Key
		rm.C = r.CopyKeys
		return
	}

	imageToMessage := func(i Image) (im ImageMessage) {
		im.K = i.Key
		im.S = i.Src
		return
	}

	pm.T = p.Title

	for _, e := range p.Elements {
		pm.E = append(pm.E, elementToMessage(e))
	}

	for _, r := range p.Repeats {
		pm.R = append(pm.R, repeatToMessage(r))
	}

	for _, i := range p.Images {
		pm.I = append(pm.I, imageToMessage(i))
	}

	return
}

func messageToPage(pm PageMessage) (p Page) {

	messageToElement := func(em ElementMessage) (e Element) {
		e.Key = em.K
		e.Text = em.T
		return
	}

	messageToRepeat := func(rm RepeatMessage) (r Repeat) {
		r.Key = rm.K
		r.CopyKeys = rm.C
		return
	}

	messageToImage := func(im ImageMessage) (i Image) {
		i.Key = im.K
		i.Src = im.S
		return
	}

	p = Page{}
	p.Title = pm.T
	for _, em := range pm.E {
		p.Elements = append(p.Elements, messageToElement(em))
	}
	for _, rm := range pm.R {
		p.Repeats = append(p.Repeats, messageToRepeat(rm))
	}
	for _, im := range pm.I {
		p.Images = append(p.Images, messageToImage(im))
	}
	return
}

func ReadPageFromJson(reader io.Reader) (*Page, error) {

	decoder := json.NewDecoder(reader)

	var pageMessage PageMessage
	err := decoder.Decode(&pageMessage)
	if err != nil {
		return nil, err
	}

	page := messageToPage(pageMessage)

	return &page, nil
}
