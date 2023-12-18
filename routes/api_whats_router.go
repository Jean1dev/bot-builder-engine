package routes

import (
	"bot_builder_engine/infra/config"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func ApiWhatsRouterHandler(w http.ResponseWriter, r *http.Request) {
	config.AllowAllOrigins(w, r)
	method := r.Method

	if method == "POST" {
		generateCode(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusBadRequest)
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	result := make([]rune, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func makeApiCall(endpoint string, method string) ([]byte, error) {
	url := fmt.Sprintf("https://whatsapp-api-da7eccbe4a89.herokuapp.com/%s", endpoint)
	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer 3n3HD2qWmeb6!hX5eQ01")

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

func generateCode(w http.ResponseWriter, r *http.Request) {
	key := generateRandomString(6)

	body, err := makeApiCall(fmt.Sprintf("instance/init?webhook=true&key=%s&token=RANDOM_STRING_HERE", key), "GET")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(string(body))

	time.Sleep(1 * time.Second)

	body, err = makeApiCall(fmt.Sprintf("instance/qr?key=%s", key), "GET")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	html := string(body)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
