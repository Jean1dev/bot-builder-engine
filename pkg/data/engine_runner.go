package data

import (
	"encoding/json"
	"fmt"

	"github.com/Jean1dev/bot-builder-engine/pkg/repository"
)

type Engine struct {
	Nodes []Node
	Edges []Edge
	Key   string
}

type EngineExecution struct {
	RequestNumber string
	NextStep      int
	Step          int
	Finished      bool
	Error         string
	Owner         string
}

func (e *Engine) GetExecutionNameForNumber(requestNumber string) string {
	return fmt.Sprintf("%s--%s.json", e.Key, requestNumber)
}

func (e *Engine) GetExecution(executionName string) (*EngineExecution, error) {
	filename := e.GetExecutionNameForNumber(executionName)
	exists := repository.VerifyIfFileExists(filename)

	if exists {
		byteValue, err := repository.Retrive(filename)
		if err != nil {
			return nil, err
		}

		var output EngineExecution
		err = json.Unmarshal(byteValue, &output)

		if err != nil {
			return nil, err
		}

		return &output, nil
	}

	output := &EngineExecution{
		RequestNumber: executionName,
		Step:          0,
		NextStep:      1,
		Finished:      false,
		Owner:         e.Key,
	}

	return output, nil
}

func (e *EngineExecution) RegisterSucess() {
	e.Step++
	e.NextStep++
	repository.SaveState(e, fmt.Sprintf("%s--%s.json", e.Owner, e.RequestNumber))
}

func (e *EngineExecution) RegisterFail(err error) {
	e.Error = err.Error()
	e.Finished = true
	repository.SaveState(e, fmt.Sprintf("%s--%s.json", e.Owner, e.RequestNumber))
}
