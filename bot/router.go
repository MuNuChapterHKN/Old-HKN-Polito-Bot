package bot

import (
	"log"
	"strings"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"github.com/marcossegovia/apiai-go"
	"gopkg.in/telegram-bot-api.v4"
)

func routeCommand(update tgbotapi.Update, ctx *BotContext) {
	command := strings.ToLower(update.Message.Command())
	var response tgbotapi.MessageConfig
	var out string
	switch command {
	case "start":
		out = "Ciao! Sono HKNBot, il bot dell'Associazione IEEE-Eta Kappa Nu" +
		" del Politecnico di Torino. Organizziamo Eventi e Gruppi di" +
		" studio, e tramite me potrai avere tutte le informazioni"+
		" di cui hai bisogno 👍  Sono un Bot testuale, per cui scrivimi e"+
		" e cercherò di risponderti!"
		response = tgbotapi.NewMessage(update.Message.Chat.ID, out)
		button1 := tgbotapi.NewKeyboardButton("Eventi \xF0\x9F\x8E\xAB")
		button2 := tgbotapi.NewKeyboardButton("Tutoraggi \xE2\x9D\x93")
		button3 := tgbotapi.NewKeyboardButton("Nascondi")
		key := tgbotapi.NewKeyboardButtonRow(button1, button2, button3)
		response.ReplyMarkup = tgbotapi.NewReplyKeyboard(key)
	case "help":
		out = "Sono un bot testuale, basato sull'"
		out += "<a href=\"https://it.wikipedia.org/wiki/Elaborazione_del_linguaggio_naturale\">NLP</a>, meglio conosciuta come "
		out += "elaborazione del linguaggio naturale.\n"
		out += "Per comunicare con me puoi scrivermi frasi e io "
		out += "cercherò di risponderti al meglio! 😊\n"
		out += "Devo ancora imparare tanto, ma per adesso"
		out += "questo è quello che puoi chiedermi\n"
		out += "Che cosa è Eta Kappa Nu?\n"
		out += "Quali sono i prossimi eventi?\n"
		out += "Avrei una domanda: ....\n"
		response = tgbotapi.NewMessage(update.Message.Chat.ID, out)
	case "domanda":
		text := update.Message.CommandArguments()
		if text != "" {
			question(update.Message.From.UserName, text, ctx.Bot)
			out = "Grazie della domanda, chiederò ai miei superiori"+
			" e ti risponderò il prima possibile! 😃"
		} else {
			out = "Per fare una domanda aggiungi il testo della domanda dopo il comando!"
		}
		response = tgbotapi.NewMessage(update.Message.Chat.ID, out)
		case "keyboard": fallthrough
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
	qr, err := ctx.ApiAi.Query(apiai.Query{Query: []string{update.Message.Text}})
	var response tgbotapi.MessageConfig
	if err != nil {
		log.Print(err)
	}
	if qr.Status.Code != 200 {
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "In questo momento"+
		" non stò molto bene, potresti provare a scrivermi più tardi?")
	}
	switch qr.Result.Action {
	case "input.question":
		question(update.Message.From.UserName, qr.Result.Params["any"], ctx.Bot)
		response = tgbotapi.NewMessage(update.Message.Chat.ID,
		qr.Result.Fulfillment.Speech)
	case "information.events":
		showEvents(update.Message.Chat.ID, ctx.Bot)
	case "buttons.hide":
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "Okay!")
		response.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	case "input.unknown":
		response = tgbotapi.NewMessage(update.Message.Chat.ID,
			qr.Result.Fulfillment.Speech +
			"Usa il comando /help per sapere cosa puoi chiedermi")
	case "input.history":
		response = tgbotapi.NewMessage(update.Message.Chat.ID, "Storia!")
	default:
		response = tgbotapi.NewMessage(update.Message.Chat.ID,
		qr.Result.Fulfillment.Speech)
	}
	if response.Text != "" {
		response.ParseMode = "HTML"
		ctx.Bot.Send(response)
	}
}

func RouteMessages(ctx *BotContext) {
	for update := range ctx.UpChannel {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			routeCommand(update, ctx)
        } else {
			routeText(update, ctx)
		}
	}
}


func showEvents(id int64, bot *tgbotapi.BotAPI) {
    msg := tgbotapi.NewMessage(id, "I prossimi eventi sono:")
    bot.Send(msg)
    // open the events file
    buf, err := ioutil.ReadFile("events.json")
    if err != nil {
        log.Print(err)
    }
    events := gjson.GetBytes(buf, "events")
    events.ForEach(func(key, value gjson.Result) bool {
        // prepare the messages from json
        out := "<b>"+value.Get("name").String()+"</b>\n"
        out += value.Get("description").String()
        msg = tgbotapi.NewMessage(id, out)
        msg.ParseMode = "HTML"
        button1 := tgbotapi.NewInlineKeyboardButtonData("Descrizione \xF0\x9F\x8E\xAB", "qu    estion")
        button2 := tgbotapi.NewInlineKeyboardButtonData("Registrati \xE2\x9D\x93", "registr    ati")
        row := tgbotapi.NewInlineKeyboardRow(button1, button2)
        keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
        msg.ReplyMarkup = keyboard
        bot.Send(msg)
        return true // keep iterating
    })
}


func question(user string, text string, bot *tgbotapi.BotAPI) {
    var out string = "<b>@"+user+":</b> "+text
    msg := tgbotapi.NewMessageToChannel("-1001111552162", out)
    msg.ParseMode = "HTML"
    bot.Send(msg)
}
