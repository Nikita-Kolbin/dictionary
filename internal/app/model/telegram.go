package model

const (
	UnknownCommandMSG = `Неизвестная команда /help`

	StartCMD = `/start`
	StartMSG = `Привет, инфа /help`

	HelpCMD = `/help`
	HelpMSG = `Список команд:

/add <word>, <translate>, <example>, <translate>
Обязательно указать слово и перевод, разделитель - запятая
пример и перевод примера оптимальные параметры

/add_time <time>
Добавить время для рассылки в формате 16:00
Не более 3 раз за день`

	AddCMD              = `/add`
	AddSuccessMSG       = `Слово "%s" добавленно!`
	AddAlreadyExistsMSG = `Это слово уже есть в коллекции`
	AddEmptyMSG         = `Слово или перевод не указано`
	AddErrorMSG         = `Не удалось добавить слово`

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
