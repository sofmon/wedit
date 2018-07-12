// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

import (
	"encoding/json"
	"sort"
)

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
