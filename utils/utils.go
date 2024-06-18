package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

func DownloadImageFromUrl(url string) (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	filename := fmt.Sprintf("file-%s", time.Microsecond.String())
	file, err := os.CreateTemp("", filename)

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return "", err
	}

	log.Printf("arquivo baixado em %s", filename)
	return file.Name(), nil
}

func InvertSlice(slice interface{}) interface{} {
	s := reflect.ValueOf(slice)
	length := s.Len()
	invertedSlice := reflect.MakeSlice(s.Type(), length, length)

	for i := 0; i < length; i++ {
		invertedSlice.Index(length - i - 1).Set(s.Index(i))
	}

	return invertedSlice.Interface()
}
