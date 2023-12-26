package application

import (
	"bot_builder_engine/data"
	"bot_builder_engine/repository"
	"bot_builder_engine/services"
	"bot_builder_engine/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
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
