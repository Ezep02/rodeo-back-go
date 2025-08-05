package sse

type SSEMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
