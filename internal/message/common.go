package msgType

type MetaData struct {
}

type Message struct {
	Payload  any      `json:"payload"`
	Metadata MetaData `json:"metadata"`
	Event    string   `json:"event"`
}

type Error struct {
	Code  int    `json:"error_code"`
	Value string `json:"error"`
}
