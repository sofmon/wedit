// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func (b *builder) findTemplatePath(path string) string {
	subFolders := strings.Split(path, "/")
	folderToCheck := b.cfg.TemplateFolder
	templateFolder := b.cfg.TemplateFolder + "/"
	for i := 0; i < len(subFolders); i++ {
		folderToCheck += "/" + subFolders[i] + "/"
		stat, err := os.Stat(folderToCheck)
		if err == nil && stat.IsDir() {
			templateFolder = folderToCheck
		} else {
			break
		}
	}
	return templateFolder
}

func (b *builder) findTemplateFile(path string) string {
	return b.findTemplatePath(path) + b.cfg.TemplateHTMLFile
}

func (b *builder) ReadPageTemplate(path string) (string, error) {
	file := b.findTemplateFile(path)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Could not read page template file '%v'. Error: %v\n", file, err)
		return "", err
	}

	return string(data), nil
}
