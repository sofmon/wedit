// Package model contains shared entities
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package model

import (
	"encoding/json"
	"sort"
)

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
