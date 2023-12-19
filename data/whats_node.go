package data

var (
	ENVIAR_MESSAGE = "ENVIAR_MENSAGEM"
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Data struct {
	Label  string `json:"label"`
	Phone  string `json:"phone,omitempty"`
	Action Action `json:"action"`
}

type DataAction struct {
	Message string `json:"message"`
}

type Action struct {
	Type string     `json:"type"`
	Data DataAction `json:"data"`
}

type Node struct {
	ID       string   `json:"id"`
	Type     string   `json:"type,omitempty"`
	Position Position `json:"position"`
	Data     Data     `json:"data"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
}

type Edge struct {
	Animated     bool   `json:"animated"`
	Type         string `json:"type"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle"`
	Target       string `json:"target"`
	TargetHandle string `json:"targetHandle,omitempty"`
	ID           string `json:"id"`
}
