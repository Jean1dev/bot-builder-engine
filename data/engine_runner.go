package data

import (
	"bot_builder_engine/repository"
	"bot_builder_engine/services"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
)

type Engine struct {
	NextStep int
	Step     int
	Nodes    []Node
	Edges    []Edge
	number   string
	Finished bool
	Key      string
	Error    error
}

func AddRunner(nodes []Node, edges []Edge, key string) {
	exists := VerifyIfKeyRegistered(key)
	if exists {
		return
	}

	engine := &Engine{
		Finished: false,
		Edges:    edges,
		Nodes:    nodes,
		Step:     0,
		NextStep: 1,
		Key:      key,
	}

	repository.SaveState(engine, fmt.Sprintf("%s.json", key))
}

func (e *Engine) registerSucess() {
	e.Step++
	e.NextStep++
	repository.SaveState(e, fmt.Sprintf("%s.json", e.Key))
}

func (e *Engine) registerFail(err error) {
	e.Error = err
	e.Finished = true
	repository.SaveState(e, fmt.Sprintf("%s.json", e.Key))
}

func VerifyIfKeyRegistered(key string) bool {
	return repository.VerifyIfFileExists(fmt.Sprintf("%s.json", key))
}

func GetEngine(key string) (*Engine, error) {
	byteValue, err := repository.Retrive(fmt.Sprintf("%s.json", key))

	if err != nil {
		log.Print(err)
		return nil, errors.New("Engine not found for suplied key")
	}

	var data Engine
	err = json.Unmarshal(byteValue, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func RunNextStep(key string, id string) {
	engine, err := GetEngine(key)
	if err != nil {
		log.Print(err)
		return
	}

	if engine.Step >= len(engine.Nodes) {
		engine.registerFail(errors.New("Engine already finished"))
		return
	}

	node := engine.Nodes[engine.Step]

	if node.Data.Action.Type == ENVIAR_MESSAGE {
		go enviarMessage(engine, &node, id)
	}
}

func enviarMessage(engine *Engine, node *Node, id string) {
	formData := url.Values{}
	formData.Set("id", id)
	formData.Set("message", node.Data.Action.Data.Message)

	if err := services.PostMessage(engine.Key, formData.Encode()); err != nil {
		engine.registerFail(err)
		return
	}

	engine.registerSucess()
}
