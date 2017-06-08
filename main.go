package main

import (
	"github.com/hkn-polito/hknbot/bot"
	"github.com/hkn-polito/hknbot/common"
)

func main() {
	configuration := common.LoadConfiguration()
	go common.WebServer()
	bot.LoadBot(configuration)
}
