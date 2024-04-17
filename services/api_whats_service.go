package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	token = func() string {
		value := os.Getenv("TOKEN")
		if value == "" {
			return "Bearer RANDOM_STRING_HERE"
		}
		return value
	}()
	baseURL = func() string {
		value := os.Getenv("BASE_URL")
		if value == "" {
			return "http://localhost:3333/"
		}
		return value
	}()
)

type TextMessageOutput struct {
	Error bool `json:"error"`
	Data  struct {
		Key struct {
			RemoteJid string `json:"remoteJid"`
			FromMe    bool   `json:"fromMe"`
			ID        string `json:"id"`
		} `json:"key"`
		Message struct {
			ExtendedTextMessage struct {
				Text string `json:"text"`
			} `json:"extendedTextMessage"`
		} `json:"message"`
		MessageTimestamp string `json:"messageTimestamp"`
		Status           string `json:"status"`
	} `json:"data"`
}

type AuditMessagesOuput struct {
	Error bool `json:"error"`
	Data  []struct {
		DbRef     string    `json:"_id"`
		CreatedAt time.Time `json:"createdAt"`
		Key       struct {
			RemoteJid string `json:"remoteJid"`
			FromMe    bool   `json:"fromMe"`
			ID        string `json:"id"`
		} `json:"key"`
		RemoteJid     string `json:"remoteJid"`
		Identificator string `json:"identificator"`
		ID            string `json:"id"`
		Messag        struct {
			ExtendedTextMessage struct {
				Text string `json:"text"`
			} `json:"extendedTextMessage"`
		} `json:"messag"`
		Status           string `json:"status"`
		MessageTimestamp struct {
			Low      int  `json:"low"`
			High     int  `json:"high"`
			Unsigned bool `json:"unsigned"`
		} `json:"messageTimestamp"`
		UpdateAt time.Time `json:"updateAt"`
	} `json:"data"`
}

type ButtonsTemplate struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

type BtnDataInput struct {
	Text    string            `json:"text"`
	Footer  string            `json:"footerText"`
	Buttons []ButtonsTemplate `json:"buttons"`
}

type InputTemplateButtonMessage struct {
	Number string       `json:"id"`
	Data   BtnDataInput `json:"btndata"`
}

func GetAuditMessages(id string) (*AuditMessagesOuput, error) {
	apiURL := fmt.Sprintf("%saudit/find?id=%s", baseURL, id)

	req, err := http.NewRequest("GET", apiURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	client := &http.Client{Timeout: time.Second * 5}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var auditMessages AuditMessagesOuput
	err = json.Unmarshal(body, &auditMessages)
	if err != nil {
		return nil, err
	}

	if auditMessages.Error {
		return nil, errors.New("Requisicao feita com erro")
	}

	return &auditMessages, nil
}

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

func PostButtonTemplate(template InputTemplateButtonMessage, key string) error {
	apiURL := fmt.Sprintf("%smessage/button?key=%s", baseURL, key)

	inputJson, err := json.Marshal(template)

	if err != nil {
		log.Println("Error marshal InputTemplateButtonMessage:", err)
		return err
	}

	payload := bytes.NewBuffer(inputJson)

	req, err := http.NewRequest("POST", apiURL, payload)

	if err != nil {
		log.Println("Erro ao criar a requisição:", err)
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

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

func sendSimpleTextMessage(key string, formData string) (*http.Response, error) {
	apiURL := fmt.Sprintf("%smessage/text?key=%s", baseURL, key)
	payload := strings.NewReader(formData)

	req, err := http.NewRequest("POST", apiURL, payload)

	if err != nil {
		log.Println("Erro ao criar a requisição:", err)
		return nil, err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Erro ao enviar a requisição:", err)
		return nil, err
	}

	return resp, nil
}

func PostMessageAndReturn(key string, formData string) (*TextMessageOutput, error) {
	resp, err := sendSimpleTextMessage(key, formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output TextMessageOutput
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	if output.Error {
		return nil, errors.New("Requisicao feita com erro")
	}

	return &output, nil

}

func PostMessage(key string, formData string) error {
	resp, err := sendSimpleTextMessage(key, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
	if res.StatusCode != 201 && res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code %s", res.Status))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}
