package model

const (
	GoodButton = `Good ✅`
	BadButton  = `Bad ❌`

	UnknownCommandMSG = `Неизвестная команда /help`

	StartCMD = `/start`
	StartMSG = `Привет, инфа /help`

	HelpCMD = `/help`
	HelpMSG = `Список команд:

/add <word>, <translate>, <example>, <translate>
Обязательно указать слово и перевод, разделитель - запятая
пример и перевод примера оптимальные параметры

/get
Возвращает одно слово на отгад

/add_time <time>
Добавить время для рассылки в формате 16:00
Не более 3 раз за день`

	AddCMD              = `/add`
	AddSuccessMSG       = `Слово "%s" добавленно!`
	AddAlreadyExistsMSG = `Это слово уже есть в коллекции`
	AddEmptyMSG         = `Слово или перевод не указано`
	AddErrorMSG         = `Не удалось добавить слово`

	GetCMD                        = `/get`
	GetSuccessWordMSG             = `Слово: %s`
	GetSuccessTranslateMSG        = `Перевод: ||%s||`
	GetSuccessExampleMSG          = `Пример: ||%s||`
	GetSuccessExampleTranslateMSG = `Перевод: ||%s||`
	GetErrorMSG                   = `Не удалось получить слово`
	GetUserHaveNotWordsMSG        = `У вас нет сохраненных слов`

	AddTimeCMD              = `/add_time`
	AddTimeSuccessMSG       = `Добавленно время для рассылки - %s`
	AddTimeAlreadyExistsMSG = `Это время и так указано`
	AddTimeEmptyMSG         = `Время не указано или указано не верно`
	AddTimeLimitMSG         = `Уже указано максимальное кол-во рассылок, сначала удалите одну через /del`
	AddTimeErrorMSG         = `Не удалось добавить время`
)

type UpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}

type Update struct {
	ID            int            `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
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

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type CallbackQuery struct {
	From            From    `json:"from"`
	Message         Message `json:"message"`
	Data            string  `json:"data"`
	InlineMessageID string  `json:"inline_message_id"`
}

type CallbackData struct {
	MessageID int  `json:"mid"`
	WordID    int  `json:"wid"`
	ChatID    int  `json:"cid"`
	Correct   bool `json:"c"`
}

type Response struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}

type Result struct {
	MessageId int `json:"message_id"`
}
