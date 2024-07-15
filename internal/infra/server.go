package infra

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jean1dev/bot-builder-engine/api"
)

func ConfigAndRunHttpServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	setupAPI()
	log.Printf("server running on %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}

}

func setupAPI() {
	http.HandleFunc("/poc/whats/generate-code", api.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/engine-run", api.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/engine-hook", api.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/verify-number", api.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/playground-send", api.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/audit", api.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/add-typebot", api.ApiWhatsRouterHandler)

	http.HandleFunc("/poc/whats/batch-send", api.BatchSend)
	http.HandleFunc("/poc/whats/batch-retrieve", api.BatchRetrive)
}
