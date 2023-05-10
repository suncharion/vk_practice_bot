package main

import (
	"fmt"
	"go_bot_tg/bot"
	"log"
	"net/http"
)

func main() {
	tgbot := bot.NewBot()
	_, err := tgbot.GetMe()
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/tgwebhook1", tgbot.UpdateWebhook)
	err = http.ListenAndServeTLS(":443", "YOURPUBLIC.pem", "YOURPRIVATE.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// основной бесконечный цикл приложения
	/*	for {


		time.Sleep(time.Second * 1)
	}*/
}
