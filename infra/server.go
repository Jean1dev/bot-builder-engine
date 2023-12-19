package main

import (
	"bot_builder_engine/routes"
	"fmt"
	"log"
	"net/http"
	"os"
)

/*
	futuramente se precisar de um roteamento mais complexo https://github.com/gorilla/mux
*/

func main() {
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
	http.HandleFunc("/poc/whats/generate-code", routes.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/engine-run", routes.ApiWhatsRouterHandler)
	http.HandleFunc("/poc/whats/engine-hook", routes.ApiWhatsRouterHandler)
}
