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

func OrchestreSending(
	nodes []data.Node,
	edges []data.Edge,
	recepients []string,
	key string) {

	steps := len(nodes)
	if steps == 0 {
		log.Print("nodes with 0 length, not processing")
		return
	}

	// sending with engine runner architecture - deprecated because is too complex
	// for _, recipient := range recepients {
	// 	AddRunner(nodes, edges, key)
	// 	go RunNextStep(key, recipient)
	// }

	actionType := nodes[0].Data.Action.Type
	if actionType == data.ENVIAR_MESSAGE {
		go sendSimpleTextMessage(recepients, key, nodes[0].Data.Action.Data.Message)
	}

	go func() {
		time.Sleep(time.Second * 5)
		getStatusSending(recepients, key, 4)
	}()
}

func sendSimpleTextMessage(recepients []string, key string, message string) {
	formData := url.Values{}
	formData.Set("message", message)

	for _, recipient := range recepients {
		formData.Set("id", recipient)
		go func() {
			result, err := services.PostMessageAndReturn(key, formData.Encode())

			if err != nil {
				log.Printf("error sending message: %v", err)
				repository.SaveState("erro na execucao", getNameForExecution(key, ""))
				return
			}

			log.Printf("result: %v", result)
		}()
	}

}

func getStatusSending(recepients []string, key string, tentatives int16) {
	result, err := services.GetAuditMessages(key)
	if err != nil {
		if tentatives > 0 {
			time.Sleep(time.Second * 5)
			getStatusSending(recepients, key, tentatives-1)
		}

		log.Printf("error getting audit messages: %v", err)
		repository.SaveState("erro na execucao", getNameForExecution(key, ""))
		return
	}

	repository.SaveState(result.Data, getNameForExecution(key, "json"))
}

func getNameForExecution(key string, format string) string {
	if format != "" {
		return fmt.Sprintf("%s_%s.%s", key, time.Now(), format)
	}

	return fmt.Sprintf("%s_%s", key, time.Now().Format("20060102150405"))
}
