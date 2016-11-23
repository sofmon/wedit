// Copyright (c) 2016, Haralampi Staykov. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
package generator

// CMS-able page
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

func NewEmptyPage() *Page {
	return &Page{
		Title:    "",
		Repeats:  []Repeat{},
		Elements: []Element{},
		Images:   []Image{},
	}
}
