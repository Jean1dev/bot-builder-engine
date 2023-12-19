package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

func SaveState(data interface{}, fileName string) error {
	file, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)

	if err != nil {
		return err
	}

	return nil
}

func Retrive(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

func VerifyIfFileExists(fileName string) bool {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return false
	}

	for _, file := range files {
		if strings.EqualFold(file.Name(), fileName) {
			return true
		}
	}

	return false
}
