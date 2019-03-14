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

	err = copyDir(cfg.TemplateFolder, cfg.PublicFolder, false)
	if err != nil {
		log.Printf("✘ unable to rebuild public folder due to an error %v\n", err)
		return
	}

	err = filepath.Walk(
		cfg.ContentFolder,
		func(path string, info os.FileInfo, theErr error) error {

			if theErr != nil {
				return theErr
			}

			if info.IsDir() {
				return nil
			}

			if !isHTML(strings.TrimSuffix(info.Name(), ".json")) {
				return nil
			}

			folderPath := filepath.Dir(path)

			relPath, err := filepath.Rel(cfg.ContentFolder, folderPath)
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
			if !strings.HasPrefix(relPath, "/") {
				relPath = "/" + relPath
			}

			paths = append(paths, relPath)

			return nil
		},
	)
	if err != nil {
		log.Printf("✘ unable to collect content due to an error: %v\n", err)
		return
	}

	sort.Strings(paths)

	for _, relPath := range paths {

		if strings.HasSuffix(relPath, "/") {
			relPath += cfg.DefaultPage
		}

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
