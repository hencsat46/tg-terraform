package models

type RequestBodyDTO struct {
	MessageInfo Message       `json:"message"`
	Callback    CallbackQuery `json:"callback_query"`
}

type Message struct {
	MessageId   int64   `json:"message_id"`
	Text        string  `json:"text"`
	ChatInfo    Chat    `json:"chat"`
	StickerInfo Sticker `json:"sticker"`
	PhotosInfo  []Photo `json:"photo"`
}

type Chat struct {
	ChatId int64 `json:"id"`
}

type Sticker struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	StickerType  string `json:"type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	IsAnimated   bool   `json:"is_animated"`
	IsVideo      bool   `json:"is_video"`
}

type Photo struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}
