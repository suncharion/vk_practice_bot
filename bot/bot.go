package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	apiUrl = "https://api.telegram.org/bot"
	apiKey = "6093483019:AAEWgJ7f4ASs1vSSmurVaJXYo-2eg-ChWRU"
)

type Bot struct {
	lastUpdate int
}

func NewBot() *Bot {
	return &Bot{}
}

type TgKeyboardButton struct {
	Text         string `json:"text"`
	Url          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

// types of data for updates

type GetUpdatesAnswer struct {
	Ok     bool
	Result []TelegramUpdate
}

type TelegramUpdate struct {
	UpdateId      int             `json:"update_id"`
	Message       TelegramMessage `json:"message"`
	CallbackQuery CallbackQuery   `json:"callback_query"`
}

type TelegramMessage struct {
	MessageId int          `json:"message_id"`
	Text      string       `json:"text"`
	Chat      TelegramChat `json:"chat"`
	From      TelegramUser `json:"from"`
}

type TelegramUser struct {
	Id int `json:"id"`
}

type TelegramChat struct {
	Id int `json:"id"`
}

type MessageAnswer struct {
	Ok bool
}

type DeleteMessageAnswer struct {
	Ok bool
}

type CallbackQuery struct {
	Id      string
	From    TelegramUser
	Message TelegramMessage
	Data    string
}

func (b *Bot) Query(method string, methodtype string, data map[string]interface{}) (string, error) {
	var resultRaw *http.Response
	var err error

	dataJSON, _ := json.Marshal(data)
	dataReader := bytes.NewBuffer(dataJSON)
	if methodtype == "GET" {
		resultRaw, err = http.Get(apiUrl + apiKey + "/" + method)
	} else {
		resultRaw, err = http.Post(apiUrl+apiKey+"/"+method, "application/json", dataReader)
	}

	//		resp, err = http.Post(endpoint, "application/json", dataReader)

	if err != nil {
		return "", err
	}
	result, _ := io.ReadAll(resultRaw.Body)
	return string(result), nil
}
func (a *Bot) GetMe() (string, error) {
	result, _ := a.Query("getMe", "GET", map[string]interface{}{})
	return result, nil
}

func (a *Bot) GetUpdates() (*GetUpdatesAnswer, error) {
	result, _ := a.Query("getUpdates?offset="+strconv.Itoa(a.lastUpdate), "GET", map[string]interface{}{})
	var parsed GetUpdatesAnswer
	fmt.Printf("%+v\n\n", result)
	err := json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		return nil, err
	}
	if len(parsed.Result) > 0 {
		a.lastUpdate = parsed.Result[len(parsed.Result)-1].UpdateId + 1
	}
	return &parsed, nil
}

// удаление inline-клавиатуры из сообщения
func (b *Bot) DeleteMessage(chat_id, message_id int) error {
	data := map[string]interface{}{
		"chat_id":    chat_id,
		"message_id": message_id,
	}
	result, err := b.Query("deleteMessage", "POST", data)
	if err != nil {
		return err
	}
	var parsed DeleteMessageAnswer
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		return err
	}
	if !parsed.Ok {
		return fmt.Errorf("Error during clearing previous message keyboard", result)
	}
	return nil
}

func (b *Bot) UpdateWebhook(w http.ResponseWriter, req *http.Request) {
	fmt.Println("got updates")

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading webhook data", err)
	}

	var update TelegramUpdate
	err = json.Unmarshal([]byte(data), &update)
	if err != nil {
		fmt.Println("Cannot parse JSON in update")
		return
	}

	if update.CallbackQuery.Message.MessageId > 0 {
		err = b.DeleteMessage(update.CallbackQuery.Message.Chat.Id, update.CallbackQuery.Message.MessageId)
		if err != nil {
			fmt.Println("Error deleting message", err)
		}
	}

	if update.CallbackQuery.Data == "pos4" {
		var keyboard []*TgKeyboardButton
		keyboard = append(keyboard, b.NewKeyboard("Ваша реклама", "https://vk.com", ""))
		keyboard = append(keyboard, b.NewKeyboard("Ваша реклама", "https://vk.com", ""))
		keyboard = append(keyboard, b.NewKeyboard("назад", "", "/start"))
		b.SendMessage(update.CallbackQuery.Message.Chat.Id, "Ваша реклама", keyboard)
	} else if update.CallbackQuery.Data == "1" {
		var keyboard []*TgKeyboardButton
		keyboard = append(keyboard, b.NewKeyboard("Twitch", "https://twitch.tv", ""))
		keyboard = append(keyboard, b.NewKeyboard("Youtube", "https://youtube.com", ""))
		keyboard = append(keyboard, b.NewKeyboard("Назад", "", "/start"))
		b.SendMessage(update.CallbackQuery.Message.Chat.Id, "Content", keyboard)
	} else if update.CallbackQuery.Data == "2" {
		var keyboard []*TgKeyboardButton
		keyboard = append(keyboard, b.NewKeyboard("Tinkoff", "https://tinkoff.ru", ""))
		keyboard = append(keyboard, b.NewKeyboard("AlfaBank", "https://alfabank.ru", ""))
		keyboard = append(keyboard, b.NewKeyboard("Назад", "", "/start"))
		b.SendMessage(update.CallbackQuery.Message.Chat.Id, "Банки", keyboard)
	} else if update.CallbackQuery.Data == "3" {
		var keyboard []*TgKeyboardButton
		keyboard = append(keyboard, b.NewKeyboard("Github", "https://github.com/suncharion", ""))
		keyboard = append(keyboard, b.NewKeyboard("Leetcode", "https://leetcode.com/", ""))
		keyboard = append(keyboard, b.NewKeyboard("Назад", "", "/start"))
		b.SendMessage(update.CallbackQuery.Message.Chat.Id, "Practice", keyboard)
	} else if update.Message.Text == "/start" || update.CallbackQuery.Data == "/start" {
		var chatId int
		if update.Message.Text != "" {
			chatId = update.Message.Chat.Id
		} else {
			chatId = update.CallbackQuery.Message.Chat.Id
		}
		var keyboard []*TgKeyboardButton

		keyboard = append(keyboard, b.NewKeyboard("Content", "", "1"))
		keyboard = append(keyboard, b.NewKeyboard("Банки", "", "2"))
		keyboard = append(keyboard, b.NewKeyboard("Practice", "", "3"))
		keyboard = append(keyboard, b.NewKeyboard("Ваша реклама", "", "pos4"))

		_, err := b.SendMessage(chatId, "VK practice bot #1", keyboard)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		b.SendMessage(update.Message.Chat.Id, "Введите /start для запуска бота", nil)
	}
}

// отправка сообщений
func (c *Bot) SendMessage(chat_id int, text string, keyboard []*TgKeyboardButton) (*MessageAnswer, error) {
	data := map[string]interface{}{
		"chat_id": chat_id,
		"text":    text,
	}
	if keyboard != nil {
		data["reply_markup"] = map[string]interface{}{
			"inline_keyboard": [][]*TgKeyboardButton{
				keyboard,
			},
		}
	}

	result, err := c.Query("sendMessage", "POST", data)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	var parsed MessageAnswer
	err = json.Unmarshal([]byte(result), &parsed)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func (d *Bot) NewKeyboard(text, url, cbdata string) *TgKeyboardButton {
	return &TgKeyboardButton{
		Text:         text,
		Url:          url,
		CallbackData: cbdata,
	}
}
