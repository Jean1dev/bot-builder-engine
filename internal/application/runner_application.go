package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/Jean1dev/bot-builder-engine/internal/services"
	"github.com/Jean1dev/bot-builder-engine/pkg/data"
	"github.com/Jean1dev/bot-builder-engine/pkg/repository"
	"github.com/Jean1dev/bot-builder-engine/utils"
)

func verifyIfKeyRegistered(key string) bool {
	return repository.VerifyIfFileExists(fmt.Sprintf("%s.json", key))
}

func getEngine(key string) (*data.Engine, error) {
	byteValue, err := repository.Retrive(fmt.Sprintf("%s.json", key))

	if err != nil {
		log.Print(err)
		return nil, errors.New("Engine not found for suplied key")
	}

	var data data.Engine
	err = json.Unmarshal(byteValue, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func AddRunner(nodes []data.Node, edges []data.Edge, key string) {
	exists := verifyIfKeyRegistered(key)
	if exists {
		return
	}

	engine := &data.Engine{
		Edges: edges,
		Nodes: nodes,
		Key:   key,
	}

	repository.SaveState(engine, fmt.Sprintf("%s.json", key))
}

func RunNextStep(key string, requestPhoneNumber string) {
	engine, err := getEngine(key)
	if err != nil {
		log.Print(err)
		return
	}

	runner, err := engine.GetExecution(requestPhoneNumber)
	if err != nil {
		log.Print(err)
		return
	}

	if runner.Step >= len(engine.Nodes) {
		runner.RegisterFail(errors.New("Engine already finished"))
		return
	}

	node := engine.Nodes[runner.Step]

	if node.Data.Action.Type == data.ENVIAR_MESSAGE {
		go enviarMessage(runner, &node, requestPhoneNumber)
	}

	if node.Data.Action.Type == data.ENVIAR_MESSAGE_IMAGEM {
		go enviarMessageWithImage(runner, &node, requestPhoneNumber)
	}

	if node.Data.Action.Type == data.ENVIAR_MESSAGE_BUTTON {
		go enviarMessageWithButton(runner, &node, requestPhoneNumber)
	}
}

func enviarMessageWithButton(engine *data.EngineExecution, node *data.Node, requestPhoneNumber string) {
	if len(node.Data.Action.Data.ActionsButtons) == 0 {
		log.Print("ActionsButtons with 0 length")
		engine.RegisterFail(errors.New("No ActionsButtons asigned"))
		return
	}

	buttonDataInput := services.BtnDataInput{
		Text:    "Titulo provisorio",
		Footer:  "Rodape provisorio",
		Buttons: make([]services.ButtonsTemplate, 0),
	}

	payloadInput := &services.InputTemplateButtonMessage{
		Data:   buttonDataInput,
		Number: requestPhoneNumber,
	}

	for _, action := range node.Data.Action.Data.ActionsButtons {
		if action.Type == data.BUTTON_TYPE_REPLY {
			template := &services.ButtonsTemplate{
				Type:    data.BUTTON_TYPE_REPLY,
				Title:   action.TitleMessage,
				Payload: action.Response,
			}

			payloadInput.Data.Buttons = append(payloadInput.Data.Buttons, *template)
		}

		if action.Type == data.BUTTON_TYPE_URL {
			template := &services.ButtonsTemplate{
				Type:    data.BUTTON_TYPE_URL,
				Title:   action.TitleMessage,
				Payload: action.Response,
			}

			payloadInput.Data.Buttons = append(payloadInput.Data.Buttons, *template)
		}

		if action.Type == data.BUTTON_TYPE_CALL {
			template := &services.ButtonsTemplate{
				Type:    data.BUTTON_TYPE_CALL,
				Title:   action.TitleMessage,
				Payload: action.Response,
			}

			payloadInput.Data.Buttons = append(payloadInput.Data.Buttons, *template)
		}
	}

	if err := services.PostButtonTemplate(*payloadInput, engine.Owner); err != nil {
		log.Print(err)
		engine.RegisterFail(err)
		return
	}

	engine.RegisterSucess()
}

func enviarMessageWithImage(engine *data.EngineExecution, node *data.Node, requestPhoneNumber string) {
	imageDownload, err := utils.DownloadImageFromUrl(node.Data.Action.Data.ImageUrl)
	if err != nil {
		log.Print(err)
		engine.RegisterFail(err)
		return
	}

	if err = services.PostMessageWithFile(engine.Owner, requestPhoneNumber, node.Data.Action.Data.ImageTitle, imageDownload); err != nil {
		log.Print(err)
		engine.RegisterFail(err)
		return
	}

	engine.RegisterSucess()
}

func enviarMessage(engine *data.EngineExecution, node *data.Node, requestPhoneNumber string) {
	formData := url.Values{}
	formData.Set("id", requestPhoneNumber)
	formData.Set("message", node.Data.Action.Data.Message)

	if err := services.PostMessage(engine.Owner, formData.Encode()); err != nil {
		engine.RegisterFail(err)
		return
	}

	engine.RegisterSucess()
}
