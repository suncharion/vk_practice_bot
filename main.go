package main

import (
	"fmt"
	"go_bot_tg/bot"
	"time"
)

func main() {
	tgbot := bot.NewBot()
	_, err := tgbot.GetMe()
	if err != nil {
		fmt.Println(err)
	}
	// основной бесконечный цикл приложения
	for {
		updates, err := tgbot.GetUpdates()
		if err != nil {
			fmt.Println(err)
		}
		// /start
		for _, update := range updates.Result {
			fmt.Printf("update: %+v\n", update)
			fmt.Println("cbdata", update.CallbackQuery.Data)

			if update.CallbackQuery.Data == "pos4" {

				var keyboard []*bot.TgKeyboardButton
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Ваша реклама", "https://vk.com", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Ваша реклама", "https://vk.com", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("назад", "", "/start"))
				tgbot.SendMessage(update.CallbackQuery.Message.Chat.Id, "Ваша реклама", keyboard)

			} else if update.CallbackQuery.Data == "1" {

				var keyboard []*bot.TgKeyboardButton
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Twitch", "https://twitch.tv", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Youtube", "https://youtube.com", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Назад", "", "/start"))
				tgbot.SendMessage(update.CallbackQuery.Message.Chat.Id, "Content", keyboard)

			} else if update.CallbackQuery.Data == "2" {

				var keyboard []*bot.TgKeyboardButton
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Tinkoff", "https://tinkoff.ru", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("AlfaBank", "https://alfabank.ru", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Назад", "", "/start"))
				tgbot.SendMessage(update.CallbackQuery.Message.Chat.Id, "Банки", keyboard)

			} else if update.CallbackQuery.Data == "3" {

				var keyboard []*bot.TgKeyboardButton
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Github", "https://github.com/suncharion", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Leetcode", "https://leetcode.com/", ""))
				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Назад", "", "/start"))
				tgbot.SendMessage(update.CallbackQuery.Message.Chat.Id, "Practice", keyboard)

			} else if update.Message.Text == "/start" || update.CallbackQuery.Data == "/start" {
				var chatId int
				if update.Message.Text != "" {
					chatId = update.Message.Chat.Id
				} else {
					chatId = update.CallbackQuery.Message.Chat.Id
				}
				var keyboard []*bot.TgKeyboardButton

				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Content", "", "1"))

				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Банки", "", "2"))

				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Practice", "", "3"))

				keyboard = append(keyboard, bot.NewBot().NewKeyboard("Ваша реклама", "", "pos4"))

				_, err := tgbot.SendMessage(chatId, "VK practice bot #1", keyboard)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				tgbot.SendMessage(update.Message.Chat.Id, "Введите /start для запуска бота", nil)
			}
			// callback_query
		}

		time.Sleep(time.Second * 1)
	}
}
