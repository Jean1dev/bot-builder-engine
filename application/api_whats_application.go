package application

import (
	"bot_builder_engine/data"
	"bot_builder_engine/services"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"
)

func GenerateQRCode(key string) (string, error) {
	// TODO :: RANDOM_STRING_HERE ???
	body, err := services.MakeApiCall(fmt.Sprintf("instance/init?webhook=true&key=%s&token=RANDOM_STRING_HERE", key), "GET")
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

	time.Sleep(1 * time.Second)

	body, err = services.MakeApiCall(fmt.Sprintf("instance/qr?key=%s", key), "GET")
	if err != nil {
		return "", err
	}

	html := string(body)
	return html, nil
}

func EngineRun(nodes []data.Node, edges []data.Edge, key string) {
	steps := len(nodes)
	if steps == 0 {
		log.Print("nodes with 0 length, not processing")
		return
	}

	AddRunner(nodes, edges, key)
}

func VerifyNumber(key string) (bool, error) {
	output, err := services.VerifyNumberOnWhatsApi(key)
	if err != nil {
		return false, err
	}

	if output.Error {
		return false, nil
	}

	if output.InstanceData.PhoneConnected {
		return true, nil
	}

	return false, nil
}

func PlayGroundSend(key string, message string, recipient string) (bool, error) {
	formData := url.Values{}
	formData.Set("message", message)
	formData.Set("id", recipient)

	output, err := services.PostMessageAndReturn(key, formData.Encode())
	if err != nil {
		return false, err
	}

	if output.Error {
		log.Printf("api whats call error")
		return false, errors.New("error sending message")
	}

	return true, nil
}
