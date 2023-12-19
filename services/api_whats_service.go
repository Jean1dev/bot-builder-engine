package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	token = "Bearer 3n3HD2qWmeb6!hX5eQ01"
)

func PostMessage(key string, formData string) error {
	apiURL := fmt.Sprintf("https://whatsapp-api-da7eccbe4a89.herokuapp.com/message/text?key=%s", key)
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
	url := fmt.Sprintf("https://whatsapp-api-da7eccbe4a89.herokuapp.com/%s", endpoint)
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}
