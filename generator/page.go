// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

// Page represents an html page
type Page struct {

	// The page title
	Title string

	// The page repeat elements
	Repeats []Repeat

	// The page variable elements
	Elements []Element

	// The page variable elements
	Images []Image
}

// NewEmptyPage creates a new empty page
func NewEmptyPage() *Page {
	return &Page{
		Title:    "",
		Repeats:  []Repeat{},
		Elements: []Element{},
		Images:   []Image{},
	}
}
