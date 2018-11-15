// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"log"
	"os"
	"strings"
)

func findTemplatePath(path string) string {
	subFolders := strings.Split(strings.Trim(path, "/"), "/")
	folderPathToCheck := cfg.TemplateFolder
	templateFolder := cfg.TemplateFolder
	for i := 0; i < len(subFolders); i++ {
		folderPathToCheck += subFolders[i] + "/"
		stat, err := os.Stat(folderPathToCheck)
		if err == nil && stat.IsDir() {
			templateFolder = folderPathToCheck
		} else {
			break
		}
	}
	return templateFolder
}

func findTemplateFile(path string) string {
	return findTemplatePath(path) + cfg.TemplateHTMLFile
}

func ReadPageTemplate(path string) (string, error) {

	data, err := prepareIncludes(path)
	if err != nil {
		log.Printf("could not read page template for path '%v'. Error: %v\n", path, err)
		return "", err
	}

	return string(data), nil
}
