package models

type Event struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    string `json:"time"`
}
