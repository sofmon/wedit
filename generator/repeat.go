// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

// Repeat represents a repeatable element on a page
type Repeat struct {
	// The unique element key in the HTML
	Key string

	// Comma separated list of keys of elements that are children of this element
	CopyKeys string
}
