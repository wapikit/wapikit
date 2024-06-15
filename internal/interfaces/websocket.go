package interfaces

type WebsocketEvent struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
