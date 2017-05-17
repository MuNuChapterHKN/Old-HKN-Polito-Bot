package main

import (
    "github.com/AntonioLangiu/hknbot/common"
    "github.com/AntonioLangiu/hknbot/bot"
)

func main() {
    configuration := common.LoadConfiguration()
    bot.LoadBot(configuration)
}
