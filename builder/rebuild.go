package builder

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func isHTML(name string) bool {
	ext := filepath.Ext(name)
	for _, e := range cfg.AllowedPageExt {
		if ext == e {
			return true
		}
	}
	return false
}

func RebuildAll() (err error) {

	var paths []string

	err = filepath.Walk(
		cfg.TemplateFolder,
		func(path string, info os.FileInfo, theErr error) error {

			if theErr != nil {
				return theErr
			}

			if info.IsDir() {
				return nil
			}

			if !isHTML(info.Name()) {
				return nil
			}

			folderPath := filepath.Dir(path)

			relPath, err := filepath.Rel(cfg.TemplateFolder, folderPath)
			if err != nil {
				log.Printf("✘ %s %v\n", path, err)
				return nil
			}

			if relPath == "." {
				relPath = ""
			}

			if !strings.HasSuffix(relPath, "/") {
				relPath += "/"
			}

			paths = append(paths, relPath)

			return nil
		},
	)

	sort.Strings(paths)

	for _, relPath := range paths {

		page, err := ReadPageData(relPath)
		if err != nil {
			log.Printf("✘ %s %v\n", relPath, err)
			continue
		}

		err = addRootData(&page)
		if err != nil {
			log.Printf("✘ %s %v\n", relPath, err)
			continue
		}

		err = WritePage(relPath, page)
		if err != nil {
			log.Printf("✘ %s %v\n", relPath, err)
			continue
		}
	}

	return
}
