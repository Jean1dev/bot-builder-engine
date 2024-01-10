package data

var (
	ENVIAR_MESSAGE        = "ENVIAR_MENSAGEM"
	ENVIAR_MESSAGE_IMAGEM = "ENVIAR_IMAGEM"
	ENVIAR_MESSAGE_BUTTON = "ENVIAR_MESSAGE_BUTTON"

	BUTTON_TYPE_REPLY = "replyButton"
	BUTTON_TYPE_URL   = "urlButton"
	BUTTON_TYPE_CALL  = "callButton"
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

type ActionButton struct {
	TitleMessage string `json:"titleMessage"`
	Type         string `json:"type"`
	Response     string `json:"response"`
}

type DataAction struct {
	Message        string         `json:"message"`
	ImageTitle     string         `json:"imageTitle"`
	ImageUrl       string         `json:"imageUrl"`
	ActionsButtons []ActionButton `json:"actionsButton"`
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
