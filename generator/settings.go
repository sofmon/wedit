package generator

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Settings caries the wedit generator configuration
type Settings struct {
	Folders struct {
		Tempalte, Data, Public string
	}
	Editor struct {
		Prefix, Host string
		Port         int
	}
}

func NewSettings(path string) (Settings, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Could not load settings file '%v'. Error: %v\n", path, err)
		return Settings{}, err
	}

	var settings Settings
	err = json.Unmarshal(file, &settings)
	if err != nil {
		log.Printf("Could not load settings from file '%v'. Error: %v\n", path, err)
		return Settings{}, err
	}

	return settings, nil
}
