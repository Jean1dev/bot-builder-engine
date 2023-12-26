package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var (
	token   = "Bearer 3n3HD2qWmeb6!hX5eQ01"
	baseURL = "https://whatsapp-api-da7eccbe4a89.herokuapp.com/"
)

func PostMessageWithFile(key string, id string, caption string, filename string) error {
	apiURL := fmt.Sprintf("%smessage/image?key=%s", baseURL, key)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	file, errFile := os.Open(filename)

	if errFile != nil {
		return errFile
	}

	defer file.Close()

	part1, errFile := writer.CreateFormFile("file", file.Name())

	_, err := io.Copy(part1, file)

	if err != nil {
		return err
	}

	_ = writer.WriteField("id", id)
	_ = writer.WriteField("caption", caption)
	err = writer.Close()

	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, payload)

	if err != nil {
		return err
	}

	req.Header.Set("Content-type", writer.FormDataContentType())
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println("Status da resposta:", res.Status)
	if res.StatusCode != 201 {
		return errors.New(fmt.Sprintf("Status code %s", res.Status))
	}

	return nil
}

func PostMessage(key string, formData string) error {
	apiURL := fmt.Sprintf("%smessage/text?key=%s", baseURL, key)
	payload := strings.NewReader(formData)

	req, err := http.NewRequest("POST", apiURL, payload)

	if err != nil {
		log.Println("Erro ao criar a requisição:", err)
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Erro ao enviar a requisição:", err)
		return err
	}

	defer resp.Body.Close()

	log.Println("Status da resposta:", resp.Status)
	if resp.StatusCode != 201 {
		return errors.New(fmt.Sprintf("Status code %s", resp.Status))
	}

	return nil
}

func MakeApiCall(endpoint string, method string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", baseURL, endpoint)
	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Authorization", token)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	log.Println("Status da resposta:", res.Status)
	if res.StatusCode != 201 {
		return nil, errors.New(fmt.Sprintf("Status code %s", res.Status))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}
