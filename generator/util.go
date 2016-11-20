// Copyright (c) 2015, Haralampi Staykov. All rights reserved. Use of this source code
// is governed by a BSD-style license that can be found in the LICENSE file.
package generator

import "bytes"

// Write set of string into bytes.Buffer
func writeBuffer(buffer *bytes.Buffer, strings ...string) error {
	for _, s := range strings {
		_, err := buffer.WriteString(s)
		if err != nil {
			return err
		}
	}
	return nil
}
