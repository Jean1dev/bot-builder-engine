package data

import (
	"bot_builder_engine/services"
	"errors"
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

type Runner struct {
	Engines []Engine
}

var (
	runner = &Runner{Engines: make([]Engine, 0)}
)

func AddRunner(nodes []Node, edges []Edge, key string) *Runner {
	exists := runner.VerifyIfKeyRegistered(key)
	if exists {
		return runner
	}

	engine := &Engine{
		Finished: false,
		Edges:    edges,
		Nodes:    nodes,
		Step:     0,
		NextStep: 1,
		Key:      key,
	}

	runner.Engines = append(runner.Engines, *engine)
	return runner
}

func Run(key string, id string) {
	runner.RunNextStep(key, id)
}

func (e *Engine) registerSucess() {
	e.Step++
	e.NextStep++
}

func (e *Engine) registerFail(err error) {
	e.Error = err
}

func (r *Runner) VerifyIfKeyRegistered(key string) bool {
	for _, engine := range r.Engines {
		if engine.Key == key {
			return true
		}
	}

	return false
}

func (r *Runner) GetEngine(key string) (*Engine, error) {
	for _, engine := range r.Engines {
		if engine.Key == key {
			return &engine, nil
		}
	}

	return nil, errors.New("Engine not found for suplied key")
}

func (r *Runner) RunNextStep(key string, id string) {
	engine, err := runner.GetEngine(key)
	if err != nil {
		log.Print(err)
		return
	}

	if engine.Step >= len(engine.Nodes) {
		engine.Finished = true
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
