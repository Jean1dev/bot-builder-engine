package routes

import (
	"bot_builder_engine/application"
	"bot_builder_engine/data"
	"bot_builder_engine/infra/config"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type EngineRegister struct {
	Node []data.Node `json:"nodes"`
	Edge []data.Edge `json:"edges"`
	Key  string      `json:"key"`
}

type IdentifierHook struct {
	RemodeJid string `json:"remoteJid"`
	FromMe    bool   `json:"fromMe"`
	ID        string `json:"id"`
}

type BodyHook struct {
	Identifier          IdentifierHook `json:"key"`
	ContactName         string         `json:"pushName"`
	VerifiedContactName string         `json:"verifiedBizName"`
	TimeStamp           int32          `json:"messageTimestamp"`
}

type Hook struct {
	Key  string   `json:"instanceKey"`
	Type string   `json:"type"`
	Body BodyHook `json:"body"`
}

func ApiWhatsRouterHandler(w http.ResponseWriter, r *http.Request) {
	config.AllowAllOrigins(w, r)
	method := r.Method

	if method == "POST" {
		if strings.HasPrefix(r.URL.Path, "/poc/whats/engine-run") {
			engine(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/poc/whats/engine-hook") {
			hook(w, r)
			return
		}

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

func hook(w http.ResponseWriter, r *http.Request) {
	var input Hook
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(input)
	data.Run(input.Key, input.Body.Identifier.RemodeJid)
}

func engine(w http.ResponseWriter, r *http.Request) {
	var input EngineRegister
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go application.EngineRun(input.Node, input.Edge, input.Key)
}

func generateCode(w http.ResponseWriter, r *http.Request) {
	key := generateRandomString(6)

	html, err := application.GenerateQRCode(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
