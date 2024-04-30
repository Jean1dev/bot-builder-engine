package application

import (
	"bot_builder_engine/data"
	"bot_builder_engine/repository"
	"bot_builder_engine/services"
	"fmt"
	"log"
	"net/url"
	"time"
)

func Retrive(id string) (string, error) {
	bytes, err := repository.Retrive(id)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func OrchestreSending(
	nodes []data.Node,
	edges []data.Edge,
	recepients []string,
	key string) string {

	steps := len(nodes)
	if steps == 0 {
		log.Print("nodes with 0 length, not processing")
		return "not processing"
	}

	// sending with engine runner architecture - deprecated because is too complex
	// for _, recipient := range recepients {
	// 	AddRunner(nodes, edges, key)
	// 	go RunNextStep(key, recipient)
	// }

	filename := getNameForExecution(key, "txt")
	actionType := nodes[0].Data.Action.Type
	if actionType == data.ENVIAR_MESSAGE {
		go sendSimpleTextMessage(recepients, key, nodes[0].Data.Action.Data.Message, filename)
	}

	return filename
}

func sendSimpleTextMessage(recepients []string, key string, message string, filename string) {
	formData := url.Values{}
	formData.Set("message", message)

	for _, recipient := range recepients {
		formData.Set("id", recipient)
		go func() {
			result, err := services.PostMessageAndReturn(key, formData.Encode())

			if err != nil {
				log.Printf("error sending message: %v", err)
				repository.SaveState("erro na execucao", filename)
				return
			}

			go repository.ApendState(fmt.Sprintf("Envio %s Id %s Status %s", recipient, result.Data.Key.ID, result.Data.Status), filename)
			go checkStatus(result.Data.Key.ID, filename, 4)
		}()
	}

}

func checkStatus(id string, filename string, tentatives int16) {
	if tentatives == 0 {
		repository.ApendState(fmt.Sprintf("ID %s Cancelando operacao pois estourou o limite de tentativas", id), filename)
		return
	}

	time.Sleep(time.Second * 2)
	result, err := services.GetAuditMessages(id)
	if err != nil {
		if tentatives > 0 {
			time.Sleep(time.Second * 1)
			checkStatus(id, filename, tentatives-1)
		}

		log.Printf("error getting audit messages: %v", err)
		repository.ApendState(fmt.Sprintf("erro ao buscar status do envio: %s %v", id, err), filename)
		return
	}

	status := result.Data[0].Status
	repository.ApendState(fmt.Sprintf("ID %s Status %s", id, status), filename)

	if status != "played" {
		checkStatus(id, filename, tentatives-1)
	}
}

func getNameForExecution(key string, format string) string {
	if format != "" {
		return fmt.Sprintf("%s_.%s", key, format)
	}

	return fmt.Sprintf("%s_", key)
}
