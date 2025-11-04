// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

import (
	"encoding/json"
	"sort"
)

// Page represents an html page
type Page struct {
	Title    string       `json:"t"`
	Repeats  PageRepeats  `json:"r"`
	Elements PageElements `json:"e"`
	Images   PageImages   `json:"i"`
}

// PageWithSettings is used to transmit page and global settings
type PageWithSettings struct {
	Page
	Settings Settings `json:"s"`
}

// NewEmptyPage creates a new empty page
func NewEmptyPage() Page {
	return Page{
		Title:    "",
		Repeats:  make(map[Key]Repeat),
		Elements: make(map[Key]Element),
		Images:   make(map[Key]Image),
	}
}

// PageElements is map[Key]Element that serializes as list
type PageElements map[Key]Element

// MarshalJSON converts a map to a json list
func (x PageElements) MarshalJSON() ([]byte, error) {
	l := []Element{}

	var keys []string
	for k := range x {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, k := range keys {
		l = append(l, x[Key(k)])
	}

	return json.Marshal(l)
}

// UnmarshalJSON converts a json list to a
func (x *PageElements) UnmarshalJSON(data []byte) (err error) {
	l := []Element{}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return
	}
	*x = make(map[Key]Element)
	for _, v := range l {
		(*x)[v.Key] = v
	}
	return
}

// PageImages is map[Key]Image that serializes as list
type PageImages map[Key]Image

// MarshalJSON converts a map to a json list
func (x PageImages) MarshalJSON() ([]byte, error) {
	l := []Image{}

	var keys []string
	for k := range x {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, k := range keys {
		l = append(l, x[Key(k)])
	}

	return json.Marshal(l)
}

// UnmarshalJSON converts a json list to a
func (x *PageImages) UnmarshalJSON(data []byte) (err error) {
	l := []Image{}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return
	}
	*x = make(map[Key]Image)
	for _, v := range l {
		(*x)[v.Key] = v
	}
	return
}

// PageRepeats is map[Key]Repeat that serializes as list
type PageRepeats map[Key]Repeat

// MarshalJSON converts a map to a json list
func (x PageRepeats) MarshalJSON() ([]byte, error) {
	l := []Repeat{}

	var keys []string
	for k := range x {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	for _, k := range keys {
		l = append(l, x[Key(k)])
	}

	return json.Marshal(l)
}

// UnmarshalJSON converts a json list to a
func (x *PageRepeats) UnmarshalJSON(data []byte) (err error) {
	l := []Repeat{}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return
	}
	*x = make(map[Key]Repeat)
	for _, v := range l {
		(*x)[v.Key] = v
	}
	return
}
