package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Jean1dev/bot-builder-engine/configs"
	"github.com/Jean1dev/bot-builder-engine/internal/application"
	"github.com/Jean1dev/bot-builder-engine/pkg/data"
)

type TypeBotInput struct {
	Key     string `json:"key"`
	ApiHost string `json:"apiHost"`
	Name    string `json:"name"`
}

type GenerateCodeInput struct {
	Code string `json:"code"`
}

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

type VerifyNumberInput struct {
	Number string `json:"number"`
	Code   string `json:"code"`
}

type PlayGroundInput struct {
	Key         string `json:"instanceKey"`
	Recipient   string `json:"recipient"`
	TextMessage string `json:"textMessage"`
}

func ApiWhatsRouterHandler(w http.ResponseWriter, r *http.Request) {
	configs.AllowAllOrigins(w, r)
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

		if strings.HasPrefix(r.URL.Path, "/poc/whats/verify-number") {
			verifyNumber(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/poc/whats/playground-send") {
			playgroundSend(w, r)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/poc/whats/add-typebot") {
			addTypeBot(w, r)
			return
		}

		generateCode(w, r)
		return
	}

	if method == "GET" {
		if strings.HasPrefix(r.URL.Path, "/poc/whats/audit") {
			listAudit(w, r)
			return
		}
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

func listAudit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("key")
	audit, err := application.ListAudit(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audit)
}

func playgroundSend(w http.ResponseWriter, r *http.Request) {
	var input PlayGroundInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := application.PlayGroundSend(input.Key, input.TextMessage, input.Recipient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func verifyNumber(w http.ResponseWriter, r *http.Request) {
	var input VerifyNumberInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	connected, err := application.VerifyNumber(input.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"connected": connected})
}

func hook(w http.ResponseWriter, r *http.Request) {
	var input Hook
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(input)
	go application.RunNextStep(input.Key, input.Body.Identifier.RemodeJid)
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
	var input GenerateCodeInput
	var key string
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		key = generateRandomString(6)
	} else {
		key = input.Code
	}

	html, err := application.GenerateQRCode(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func addTypeBot(w http.ResponseWriter, r *http.Request) {
	var input TypeBotInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success := application.AddTypeBotOnNumber(input.Key, input.ApiHost, input.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": success})
}
