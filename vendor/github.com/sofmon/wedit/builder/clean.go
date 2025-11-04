package builder

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CleanUp content and public folder from items not matching template folder
func CleanUp() (err error) {

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

			_, err = findTemplatePath(relPath)
			if err == nil {
				return nil
			}
			if err != ErrPathMismatch {
				log.Printf("✘ %s %v\n", path, err)
				return nil
			}

			// the err is ErrPathMismatch so we need to clean it up
			log.Printf("🗑 %s\n", relPath)

			err = os.Remove(path)
			if err != nil {
				log.Printf("✘ %s %v\n", path, err)
				return nil
			}

			return nil
		},
	)
	if err != nil {
		return
	}

	err = filepath.Walk(
		cfg.ContentFolder,
		func(path string, info os.FileInfo, theErr error) error {

			if theErr != nil {
				return theErr
			}

			if !info.IsDir() {
				return nil
			}

			empty, err := isFolderEmpty(path)
			if err != nil {
				log.Printf("✘ %s %v\n", path, err)
				return nil
			}

			if empty {
				err = os.Remove(path)
				if err != nil {
					log.Printf("✘ %s %v\n", path, err)
					return nil
				}

			}

			return nil
		},
	)

	return
}

func isFolderEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
