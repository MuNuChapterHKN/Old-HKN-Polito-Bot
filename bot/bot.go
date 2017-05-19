package bot

import (
    "log"
    "github.com/marcossegovia/apiai-go"
    "github.com/AntonioLangiu/hknbot/common"
    "gopkg.in/telegram-bot-api.v4"
)

type BotContext struct {
    Config *common.Configuration
    Bot *tgbotapi.BotAPI
    UpChannel <-chan tgbotapi.Update
    ApiAi *apiai.ApiClient
}

func LoadBot(configuration *common.Configuration) {
    ctx := InitBot(configuration)
	configApiAi(ctx)
    RouteMessages(ctx)
}

func InitBot(configuration *common.Configuration) *BotContext {
    ctx := BotContext{}
    ctx.Config = configuration

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramAPI)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
    ctx.Bot = bot
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	ctx.UpChannel, err = bot.GetUpdatesChan(u)
    if err != nil {
        log.Panic(err)
    }

    return &ctx
}

func configApiAi(ctx *BotContext) {
    client, err := apiai.NewClient(
		&apiai.ClientConfig{
			Token:      ctx.Config.ApiAiToken,
            QueryLang:  ctx.Config.ApiAiQueryLang,
		},
	)
	if err != nil {
        log.Panic(err)
	}
    ctx.ApiAi = client
}
