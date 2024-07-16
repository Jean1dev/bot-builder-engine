package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Jean1dev/bot-builder-engine/internal/infra"
	"github.com/Jean1dev/bot-builder-engine/internal/services"
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
	services.ConfigureEnv()
	infra.ConfigAndRunHttpServer()
}
