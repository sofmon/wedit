package generator

import "os"
import "io/ioutil"
import "encoding/json"
import "log"
import "strings"

const dataFileName = "index.json"
const templateFileName = "index.html"

type FileExplorer struct {
	settings Settings
}

func NewFileExplorer(settings Settings) FileExplorer {
	return FileExplorer{settings}
}

func (f FileExplorer) ReadPageTemplate(path string) (string, error) {
	subFolders := strings.Split(path, "/")
	folderToCheck := f.settings.Folders.Tempalte
	templateFolder := f.settings.Folders.Tempalte
	for i := 0; i < len(subFolders); i++ {
		folderToCheck += subFolders[i] + "/"
		stat, err := os.Stat(folderToCheck)
		if err == nil && stat.IsDir() {
			templateFolder = folderToCheck
		} else {
			break
		}
	}

	file := templateFolder + templateFileName
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Could not read page template file '%v'. Error: %v\n", file, err)
		return "", err
	}

	return string(data), nil
}

func (f FileExplorer) ReadPageData(path string) (Page, error) {

	file := f.settings.Folders.Data + path + dataFileName
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Could not load page data file '%v'. Error: %v\n", file, err)
		return Page{}, err
	}

	var page Page
	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("Could not load page data file '%v'. Error: %v\n", file, err)
		return Page{}, err
	}

	return page, nil
}

func (f FileExplorer) WritePageData(path string, page Page) error {

	folder := f.settings.Folders.Data + path
	os.MkdirAll(folder, 666)

	file := folder + dataFileName

	data, err := json.Marshal(page)
	if err != nil {
		log.Printf("Unable to save page data at '%v'. Error: %v", file, err)
		return err
	}

	err = ioutil.WriteFile(file, data, 666)
	if err != nil {
		log.Printf("Unable to save page data at '%v'. Error: %v", file, err)
		return err
	}

	return nil
}

func (f FileExplorer) WriteAsset(path, asset string, data []byte) error { // TODO
	folder := f.settings.Folders.Data + path
	os.MkdirAll(folder, 666)

	file := folder + asset
	err := ioutil.WriteFile(file, data, 666)
	if err != nil {
		log.Printf("Unable to save asset at '%v'. Error: %v", file, err)
		return err
	}

	return nil
}
