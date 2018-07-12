// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

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
