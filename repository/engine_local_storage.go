package repository

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

// RACE CONDITION AQUI
func ApendState(data interface{}, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return nil
}

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

	byteValue, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

func VerifyIfFileExists(fileName string) bool {
	files, err := os.ReadDir(".")
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
