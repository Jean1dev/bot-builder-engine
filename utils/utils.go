package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
