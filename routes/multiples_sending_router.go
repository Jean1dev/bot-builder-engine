package routes

import (
	"bot_builder_engine/application"
	"bot_builder_engine/infra/config"
	"encoding/json"
	"log"
	"net/http"
)

type MultipleSendingInputBody struct {
	To          []string       `json:"to"`
	Engine      EngineRegister `json:"engine"`
	KeyWhatsApp string         `json:"key"`
	ExternalId  string         `json:"external_id"`
}

func BatchSend(w http.ResponseWriter, r *http.Request) {
	config.AllowAllOrigins(w, r)
	var inputBody MultipleSendingInputBody
	if err := json.NewDecoder(r.Body).Decode(&inputBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	application.OrchestreSending(inputBody.Engine.Node, inputBody.Engine.Edge, inputBody.To, inputBody.KeyWhatsApp)

	log.Print(inputBody)
}
