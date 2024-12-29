package model

import "errors"

var ErrUnknownCommand = errors.New("unknown command")

const (
	StartCMD = "/start"
	HelpCMD  = "/help"

	UnknownCommandMSG = `Неизвестная команда /help`
	StartMSG          = `Привет, инфа /help`
	HelpMSG           = `Список команд:`
)

type UpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
