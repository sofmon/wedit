package builder

import (
	"log"
	"os"
	"path/filepath"
)

func (b *builder) RebuildAll() (err error) {

	filepath.Walk(
		b.cfg.TemplateFolder,
		func(path string, info os.FileInfo, theErr error) error {

			if theErr != nil {
				return theErr
			}

			if info.IsDir() {
				return nil
			}

			if info.Name() != b.cfg.TemplateHTMLFile {
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

			local, err := b.ReadPageData(relPath)
			if err != nil {
				log.Printf("%40s - not processed due %v\n", relPath, err)
				return nil
			}

			page, err := b.addRootData(local)
			if err != nil {
				log.Printf("%40s - not processed due %v\n", relPath, err)
				return nil
			}

			err = b.WritePage(relPath, page)
			if err != nil {
				log.Printf("%40s - not processed due %v\n", relPath, err)
				return nil
			}

			log.Printf("%40s - rebuild\n", relPath)

			return nil
		},
	)

	return
}
