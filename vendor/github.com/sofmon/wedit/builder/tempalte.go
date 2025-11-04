// Package builder builds the wedit public folder
// Part of `wedit` project (https://wedit.io) (https://github.com/sofmon/wedit)
// Use of this source code is governed by MIT license that can be found in the LICENSE file.
package builder

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrPathMismatch is used when path can not be matched with tempalte folder
	ErrPathMismatch = errors.New("provided path could not be matched to template folder ")
)

func findTemplatePath(path string) (string, error) {
	subFolders := strings.Split(strings.Trim(path, "/"), "/")
	templateFolder := cfg.TemplateFolder
	for i := 0; i < len(subFolders); i++ {
		// Check if folder exists
		templateCandidate := filepath.Join(templateFolder, subFolders[i])
		stat, err := os.Stat(templateCandidate)
		if err == nil && stat.IsDir() {
			templateFolder = templateCandidate
			continue
		}

		// check if matches a variation folder
		templateCandidate = filepath.Join(templateFolder, cfg.VariationFolderName)
		stat, err = os.Stat(templateCandidate)
		if err == nil && stat.IsDir() {
			templateFolder = templateCandidate
			continue
		}

		return "", ErrPathMismatch
	}
	return templateFolder, nil
}

func findTemplateFile(path string) (string, error) {
	dir, file := filepath.Split(path)
	templatePath, err := findTemplatePath(dir)
	if err != nil {
		return "", err
	}
	return filepath.Join(templatePath, file), nil
}

func ReadPageTemplate(path string) (string, error) {

	data, err := prepareIncludes(path)
	if err != nil {
		log.Printf("✘ could not read page template for path '%v'. Error: %v\n", path, err)
		return "", err
	}

	return string(data), nil
}
