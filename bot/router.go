package bot

import (
    "fmt"
	"github.com/marcossegovia/apiai-go"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
    "strconv"
)

func RouteMessages(ctx *BotContext) {
	for update := range ctx.UpChannel {
        // Message to channel
        if update.ChannelPost != nil {
            routeChannel(update, ctx)
            continue
        }
        // Callback query
        if update.CallbackQuery != nil {
            routeCallbackQuery(update, ctx)
        }
        // Message from Chat
        if update.Message == nil {
			continue
		}
        if update.Message.IsCommand() {
			routeCommand(update, ctx)
		} else if update.Message.Text == "" {
			routeInvalid(update, ctx)
		} else {
			routeText(update, ctx)
		}
	}
}

func routeCallbackQuery(update tgbotapi.Update, ctx *BotContext) {
    data := update.CallbackQuery.Data
    split := strings.Split(data, ":")
    switch split[0] {
    case "event":
        RouteEventQuery(split, update, ctx)
    default:
        log.Print("error, I received a callback query I'm not ready to handle")
    }
}

func routeChannel(update tgbotapi.Update, ctx *BotContext) {
    if update.ChannelPost.Chat.ID == -1001111552162 {
        if update.ChannelPost.ReplyToMessage != nil {
            if update.ChannelPost.Text == "" {
                return
            }
            text := update.ChannelPost.ReplyToMessage.Text
            split := strings.Split(text, "@")
            if len(split) > 0 {
                id, err := strconv.ParseInt(split[0], 10, 64)
                if err != nil {
                    log.Print(err)
                    return
                }
                out := "<b>Ecco la risposta che ho ricevuto dai miei capi:</b> "
                out += update.ChannelPost.Text
                response := tgbotapi.NewMessage(id, out)
                response.ParseMode = "HTML"
                ctx.Bot.Send(response)
            }
        }
    }
}

func routeCommand(update tgbotapi.Update, ctx *BotContext) {
	command := strings.ToLower(update.Message.Command())
	var response tgbotapi.MessageConfig
	var out string
	switch command {
	case "start":
		StartMessage(update, ctx)
	case "help":
		HelpMessage(update, ctx)
	case "domanda":
		text := update.Message.CommandArguments()
		if text != "" {
			question(update.Message.From.ID, update.Message.From.UserName, text, ctx.Bot)
			out = "Grazie della domanda, chiederò ai miei superiori" +
				" e ti risponderò il prima possibile! 😃"
		} else {
			out = "Per fare una domanda aggiungi il testo della domanda dopo il comando!"
		}
		response = tgbotapi.NewMessage(update.Message.Chat.ID, out)
	case "keyboard":
		fallthrough
	case "tastiera":
		out = "Ai tuoi ordini!"
		response = tgbotapi.NewMessage(update.Message.Chat.ID, out)
		button1 := tgbotapi.NewKeyboardButton("Eventi \xF0\x9F\x8E\xAB")
		button2 := tgbotapi.NewKeyboardButton("Tutoraggi \xE2\x9D\x93")
		button3 := tgbotapi.NewKeyboardButton("Nascondi")
		key := tgbotapi.NewKeyboardButtonRow(button1, button2, button3)
		response.ReplyMarkup = tgbotapi.NewReplyKeyboard(key)
	default:
		out = "Scusa, non ho capito! Prova a scrivere il testo senza lo slash"
		response = tgbotapi.NewMessage(update.Message.Chat.ID, out)
	}
	if response.Text != "" {
		response.ParseMode = "HTML"
		ctx.Bot.Send(response)
	}
}

func routeText(update tgbotapi.Update, ctx *BotContext) {
	qr, err := ctx.ApiAi.Query(apiai.Query{Query: []string{update.Message.Text}, SessionId: ctx.Config.ApiAiSessionId})
	var response tgbotapi.MessageConfig
	if err != nil {
		log.Print(err)
		return
	}
	if qr.Status.Code != 200 {
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "In questo momento"+
			" non stò molto bene, potresti provare a scrivermi più tardi?")
	}
	switch qr.Result.Action {
	case "input.question":
		question(update.Message.From.ID, update.Message.From.UserName, qr.Result.Params["any"], ctx.Bot)
		response = tgbotapi.NewMessage(update.Message.Chat.ID,
			qr.Result.Fulfillment.Speech)
	case "information.events":
		ShowEvents(update, ctx)
	case "buttons.hide":
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "Okay!")
		response.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case "input.unknown":
		response = tgbotapi.NewMessage(update.Message.Chat.ID,
			qr.Result.Fulfillment.Speech+
				"Usa il comando /help per sapere cosa puoi chiedermi")
	case "input.history":
		HknHistory(update, ctx)
	case "support.problem":
		question(update.Message.From.ID, update.Message.From.UserName, "bug_reporting: "+update.Message.Text, ctx.Bot)
	default:
        if qr.Result.Fulfillment.Messages != nil {
            for _,elem := range qr.Result.Fulfillment.Messages {
                if elem.Type == 0 {
                    response = tgbotapi.NewMessage(update.Message.Chat.ID,
                    elem.Speech)
                }
            }
        } else {
            response = tgbotapi.NewMessage(update.Message.Chat.ID,
            qr.Result.Fulfillment.Speech)
        }
	}
	if response.Text != "" {
		response.ParseMode = "HTML"
		ctx.Bot.Send(response)
	}
}

func routeInvalid(update tgbotapi.Update, ctx *BotContext) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, "Scusa ma io capisco solo il testo!")
	ctx.Bot.Send(response)
}

func question(id int, user string, text string, bot *tgbotapi.BotAPI) {
    out := fmt.Sprintf("%v", id)
    out += "<b>@" + user + ":</b> " + text
	msg := tgbotapi.NewMessageToChannel("-1001111552162", out)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}
