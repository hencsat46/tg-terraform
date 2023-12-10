package models

type ResponseBodyTextDTO struct {
	ChatID int64           `json:"chat_id"`
	Text   string          `json:"text"`
	Entity []MessageEntity `json:"entities"`
}

type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	Language string `json:"language"`
}
