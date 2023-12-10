package models

type ResponseBodyButtonsDTO struct {
	ChatID  int64                `json:"chat_id"`
	Text    string               `json:"text"`
	Buttons InlineKeyboardMarkup `json:"reply_markup"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButtons `json:"inline_keyboard"`
}

type InlineKeyboardButtons struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type CallbackQuery struct {
	ID          string  `json:"id"`
	MessageInfo Message `json:"message"`
	Data        string  `json:"data"`
}

type AnswerCallback struct {
	Id   string `json:"callback_query_id"`
	Text string `json:"text"`
}
