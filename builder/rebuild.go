package builder

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (b *builder) RebuildAll() (err error) {

	var paths []string

	err = filepath.Walk(
		b.cfg.TemplateFolder,
		func(path string, info os.FileInfo, theErr error) error {

			if theErr != nil {
				return theErr
			}

			if info.IsDir() {
				return nil
			}

			if !b.cfg.Contains(b.cfg.PageFiles, info.Name()) {
				return nil
			}

			folderPath := filepath.Dir(path)

			relPath, err := filepath.Rel(b.cfg.TemplateFolder, folderPath)
			if err != nil {
				log.Printf("%40s - not processed due %v\n", path, err)
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

		for _, pageFile := range b.cfg.PageFiles {

			page, err := b.ReadPageData(relPath, pageFile)
			if err != nil {
				log.Printf("%40s - not processed due %v\n", relPath, err)
				continue
			}

			err = b.addRootData(&page)
			if err != nil {
				log.Printf("%40s - not processed due %v\n", relPath, err)
				continue
			}

			err = b.WritePage(relPath, page, pageFile)
			if err != nil {
				log.Printf("%40s - not processed due %v\n", relPath, err)
				continue
			}
		}
	}

	return
}
