// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sofmon/wedit/model"
)

var (
	rootData     *model.Page
	rootDataLock sync.Mutex
)

func (b *builder) getRootData() (page model.Page, err error) {

	rootDataLock.Lock()
	defer rootDataLock.Unlock()

	if rootData != nil {
		page = *rootData
		return
	}

	file := b.cfg.ContentFolder + b.cfg.RootJSONFile
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			page = model.NewEmptyPage()
			err = nil // it is an empty data
			return
		}
		log.Printf("Could not load root data file '%v'. Error: %v\n", file, err)
		return
	}

	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("Could not load root data file '%v'. Error: %v\n", file, err)
		return
	}

	if page.Elements == nil {
		page.Elements = make(model.PageElements)
	}
	if page.Images == nil {
		page.Images = make(model.PageImages)
	}
	if page.Repeats == nil {
		page.Repeats = make(model.PageRepeats)
	}

	rootData = &page

	return
}

func (b *builder) clearRootData() {

	rootDataLock.Lock()
	defer rootDataLock.Unlock()

	rootData = nil
}

func (b *builder) updateRootData(page model.Page) (wasUpdated bool, err error) {

	wasUpdated = false

	newRootData, err := b.getRootData()
	if err != nil {
		return
	}

	rootDataLock.Lock()
	defer rootDataLock.Unlock()

	for k, v := range page.Elements {
		if cv, ok := newRootData.Elements[k]; !ok || !isSameElement(cv, v) {
			newRootData.Elements[k] = v
			wasUpdated = true
		}
	}

	for k, v := range page.Images {
		if cv, ok := newRootData.Images[k]; !ok || !isSameImage(cv, v) {
			newRootData.Images[k] = v
			wasUpdated = true
		}
	}

	for k, v := range page.Repeats {
		if cv, ok := newRootData.Repeats[k]; !ok || !isSameRepeat(cv, v) {
			newRootData.Repeats[k] = v
			wasUpdated = true
		}
	}

	if !wasUpdated {
		return
	}

	rootData = &newRootData

	file := b.cfg.ContentFolder + b.cfg.RootJSONFile

	data, err := json.MarshalIndent(newRootData, "", "  ")
	if err != nil {
		log.Printf("unable to save root data at '%v'; error: %v", file, err)
		return false, err
	}

	err = ioutil.WriteFile(file, data, 0777)
	if err != nil {
		log.Printf("unable to save root data at '%v'; error: %v", file, err)
		return false, err
	}

	return
}

func isSameElement(x, y model.Element) bool {

	if x.Key != y.Key {
		return false
	}

	if x.Text != y.Text {
		return false
	}

	return true
}

func isSameImage(x, y model.Image) bool {

	if x.Key != y.Key {
		return false
	}

	// TODO !!!

	return true
}

func isSameRepeat(x, y model.Repeat) bool {
	if x.Key != y.Key {
		return false
	}

	if len(x.CopyKeys) != len(y.CopyKeys) {
		return false
	}

	for i, k := range x.CopyKeys {
		if k != y.CopyKeys[i] {
			return false
		}
	}

	return true
}

func (b *builder) splitRootData(page model.Page) (local model.Page, root model.Page) {

	local.Title = page.Title
	local.Repeats = make(map[model.Key]model.Repeat)
	local.Elements = make(map[model.Key]model.Element)
	local.Images = make(map[model.Key]model.Image)

	root.Repeats = make(map[model.Key]model.Repeat)
	root.Elements = make(map[model.Key]model.Element)
	root.Images = make(map[model.Key]model.Image)

	for k, v := range page.Elements {
		if strings.HasPrefix(string(k), b.cfg.RootKeyPrefix) {
			root.Elements[k] = v
		} else {
			local.Elements[k] = v
		}
	}

	for k, v := range page.Images {
		if strings.HasPrefix(string(k), b.cfg.RootKeyPrefix) {
			root.Images[k] = v
		} else {
			local.Images[k] = v
		}
	}

	for k, v := range page.Repeats {
		if strings.HasPrefix(string(k), b.cfg.RootKeyPrefix) {
			root.Repeats[k] = v
		} else {
			local.Repeats[k] = v
		}
	}

	return
}

func (b *builder) addRootData(page *model.Page) error {

	global, err := b.getRootData()
	if err != nil {
		return err
	}

	if page.Elements == nil {
		page.Elements = make(map[model.Key]model.Element)
	}
	if page.Repeats == nil {
		page.Repeats = make(map[model.Key]model.Repeat)
	}
	if page.Images == nil {
		page.Images = make(map[model.Key]model.Image)
	}

	for k, v := range global.Elements {
		page.Elements[k] = v
	}

	for k, v := range global.Images {
		page.Images[k] = v
	}

	for k, v := range global.Repeats {
		page.Repeats[k] = v
	}

	return nil
}
