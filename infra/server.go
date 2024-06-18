package main

import (
	"bot_builder_engine/routes"
	"bot_builder_engine/services"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
	futuramente se precisar de um roteamento mais complexo https://github.com/gorilla/mux
*/

func init() {
	isHerokuEnv := os.Getenv("HEROKU_DYNO")
	if isHerokuEnv != "" {
		return
	}

	// Load environment variables from .env file"
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}
	fmt.Println("Current directory:", currentDir)

	envFilePath := currentDir + "/../.env"
	file, err := os.Open(envFilePath)
	if err != nil {
		log.Fatal("Error opening .env file:", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading .env file:", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	setupAPI()
	services.ConfigureEnv()
	log.Printf("server running on %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}

}

func setupAPI() {
	http.HandleFunc("/poc/whats/generate-code", routes.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/engine-run", routes.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/engine-hook", routes.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/verify-number", routes.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/playground-send", routes.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/audit", routes.ApiWhatsRouterHandler)

	http.HandleFunc("/poc/whats/batch-send", routes.BatchSend)
	http.HandleFunc("/poc/whats/batch-retrieve", routes.BatchRetrive)
}
