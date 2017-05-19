package main

import (
	"github.com/AntonioLangiu/hknbot/bot"
	"github.com/AntonioLangiu/hknbot/common"
)

func main() {
	configuration := common.LoadConfiguration()
	go common.WebServer()
	bot.LoadBot(configuration)
}
