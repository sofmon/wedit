// Package service handles wedit HTTP calls
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package service

// Config for the wedit http service
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
