// Copyright (c) 2016, Haralampi Staykov. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
package generator

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

func WritePageAsEscapedJson(path, host string, page *Page, buffer *bytes.Buffer) error {

	err := writeBuffer(buffer,
		"{\\\"", PAGE_SITE, "\\\":\\\"", "",
		"\\\",\\\"", PAGE_PATH, "\\\":\\\"", path,
		"\\\",\\\"", PAGE_HOST, "\\\":\\\"", host,
		"\\\",\\\"", PAGE_TITLE, "\\\":\\\"", page.Title,
		"\\\",\\\"", PAGE_ELEMENTS, "\\\":[")
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

	err = writeBuffer(buffer, "],\\\"", PAGE_REPEATS, "\\\":[")
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

	_, err = buffer.WriteString("]}")
	if err != nil {
		return err
	}

	return nil
}

func WriteElementAsEscapedJson(element *Element, buffer *bytes.Buffer) error {
	err := writeBuffer(buffer,
		"{\\\"", ELEMENT_KEY, "\\\":\\\"", element.Key,
		"\\\",\\\"", ELEMENT_TEXT, "\\\":\\\"", strings.Replace(element.Text, "'", "\\'", -1),
		"\\\",\\\"i1\\\":\\\"", "",
		"\\\",\\\"i2\\\":\\\"", "",
		"\\\"}")
	if err != nil {
		return err
	}
	return nil
}

func WriteRepeatAsEscapedJson(repeat *Repeat, buffer *bytes.Buffer) error {
	err := writeBuffer(buffer,
		"{\\\"", REPEAT_KEY, "\\\":\\\"", repeat.Key,
		"\\\",\\\"", REPEAT_COPY_KEYS, "\\\":\\\"", repeat.CopyKeys,
		"\\\"}")
	if err != nil {
		return err
	}

	return nil
}

type PageMessage struct {
	T    string           // Title
	E    []ElementMessage // Elements
	R    []RepeatMessage  // Repeats
	Html string           // HTML
}

type ElementMessage struct {
	K string // Key
	T string // Text
}

type RepeatMessage struct {
	K string // Key
	C string // CopyKeys
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

	pm.T = p.Title

	for _, e := range p.Elements {
		pm.E = append(pm.E, elementToMessage(e))
	}

	for _, r := range p.Repeats {
		pm.R = append(pm.R, repeatToMessage(r))
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

	p = Page{}
	p.Title = pm.T
	p.HTML = pm.Html
	for _, em := range pm.E {
		p.Elements = append(p.Elements, messageToElement(em))
	}
	for _, rm := range pm.R {
		p.Repeats = append(p.Repeats, messageToRepeat(rm))
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
