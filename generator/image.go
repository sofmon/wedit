// Copyright (c) 2016, Haralampi Staykov (http://haralampi.com). All rights reserved.
// Use of this source code is governed by MIT license that can be found in the LICENSE file.

package generator

// Image represents an image on a page
type Image struct {

	// The unique element key in the HTML
	Key string

	// The serving src for the image
	Src string
}
