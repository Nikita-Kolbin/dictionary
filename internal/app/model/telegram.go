package model

const (
	GoodButton = `Good ✅`
	BadButton  = `Bad ❌`

	UnknownCommandMSG = `Неизвестная команда, постморите /help`

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
Добавить время для рассылки в формате 16:00 (по МСК)
Не более 3 раз за день

/get_time
Возвращает время на которое у вас стоит рассылка

/del_time <time>
Удаляет время рассылки

/set_count <count>
Установить сколько слов будет в рассылке
count - целое число от 1 до 25`

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
	GetSuccessOpenInTranslator    = `[Переводчик](%s)`
	GetErrorMSG                   = `Не удалось получить слово`
	GetUserHaveNotWordsMSG        = `У вас нет сохраненных слов`

	AddTimeCMD              = `/add_time`
	AddTimeSuccessMSG       = `Добавленно время для рассылки - %s`
	AddTimeAlreadyExistsMSG = `Это время и так указано`
	AddTimeEmptyMSG         = `Время не указано или указано неверно`
	AddTimeLimitMSG         = `Уже указано максимальное кол-во рассылок, сначала удалите одну через /del_time`
	AddTimeErrorMSG         = `Не удалось добавить время`

	GetTimeCMD        = `/get_time`
	GetTimeSuccessMSG = `Время рассылки (по МСК):`
	GetTimeEmptyMSG   = `У вас не установленно время рассылки`
	GetTimeErrorMSG   = `Ошибка получения времени рассылки`

	DelTimeCMD        = `/del_time`
	DelTimeSuccessMSG = `Время рассылки %s удалено`
	DelTimeEmptyMSG   = `У вас и так нет рассылки на это время`
	DelTimeErrorMSG   = `Ошибка удаления времени рассылки`

	SetCountCMD             = `/set_count`
	SetCountSuccessMSG      = `Установленно кол-во слов для рассылки: %d`
	SetCountEmptyMSG        = `Кол-во не указано или указано неверно`
	SetCountUserNotFoundMSG = `Пользователь не найден, используйте /start для обновления базы`
	SetCountErrorMSG        = `Ошибка установки кол-ва слов для рассылки`
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
	MessageID int `json:"message_id"`
}
