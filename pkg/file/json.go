package file

import (
	"encoding/json"
	"io/ioutil"
)


func LoadStructuresFromDir(path string, structure interface{}) error {


	return nil
}

func LoadStructureFromJsonFile(path string, structure interface{}) error {
	if FileNotExists(path) {
		err := ReaderError{Message: "File does not exists on path: " + path}
		return &err
	}
	configFileData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	parseError := json.Unmarshal(configFileData, structure)
	if parseError != nil {
		return parseError
	}

	return nil
}


