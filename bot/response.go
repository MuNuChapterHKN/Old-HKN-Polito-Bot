package bot

import (
	"gopkg.in/telegram-bot-api.v4"
)

func SendText(update tgbotapi.Update, ctx *BotContext, text string) {
	response := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	response.ParseMode = "HTML"
	ctx.Bot.Send(response)
}

func StartMessage(update tgbotapi.Update, ctx *BotContext) {
	out := "Ciao! Sono HKNBot, il bot dell'Associazione IEEE-Eta Kappa Nu" +
		" del Politecnico di Torino. Organizziamo Eventi e Gruppi di" +
		" studio, e tramite me potrai avere tutte le informazioni" +
		" di cui hai bisogno 👍  Sono un Bot testuale, per cui scrivimi e" +
		" e cercherò di risponderti!"
	response := tgbotapi.NewMessage(update.Message.Chat.ID, out)
	button1 := tgbotapi.NewKeyboardButton("Eventi \xF0\x9F\x8E\xAB")
	button2 := tgbotapi.NewKeyboardButton("Tutoraggi \xE2\x9D\x93")
	button3 := tgbotapi.NewKeyboardButton("Nascondi")
	key := tgbotapi.NewKeyboardButtonRow(button1, button2, button3)
	response.ReplyMarkup = tgbotapi.NewReplyKeyboard(key)
	response.ParseMode = "HTML"
	ctx.Bot.Send(response)
	sendCommands(update, ctx)
}

func HelpMessage(update tgbotapi.Update, ctx *BotContext) {
	out := "Sono un bot testuale, basato sull'"
	out += "NLP, meglio conosciuta come "
	out += "elaborazione del linguaggio naturale.\n"
	out += "Per comunicare con me puoi scrivermi frasi e io "
	out += "cercherò di risponderti al meglio! 😊\n"
	SendText(update, ctx, out)
	sendCommands(update, ctx)
}

func sendCommands(update tgbotapi.Update, ctx *BotContext) {
    out := "Queste sono <b>alcune delle cose che puoi chiedermi:</b>\n"
	out += "  Che cosa è Eta Kappa Nu?\n"
	out += "  Quali sono i prossimi eventi?\n"
	out += "  Qual'è la storia di Eta Kappa Nu?\n\n"
	out += "Sono inoltre capace di ripondere ad altre domande, "
	out += "prova a scoprire quali! 😀\n\n"

	out += "Sono inoltre capace di inoltrare delle domande al team di Eta Kappa Nu, "
	out += "per fare questo puoi usare il comando /domanda "
	out += "e aggiungere il testo dopo, o scrivere \n"
	out += "\"Avrei una domanda:\" seguito dal testo della domanda stessa"
	SendText(update, ctx, out)
}

func HknHistory(update tgbotapi.Update, ctx *BotContext) {
	out := "IEEE-Eta Kappa Nu è una honor society "
	out += "fondata in Ottobre 1904 da Maurice L.Carr nell'università dell Illinois.\n "
	out += "Il suo scopo è riunire i migliori studenti delle facoltà di Ingegneria Informatica "
	out += "ed elettronica e conta oggi più di 200,000 membri in tutto il mondo!\n"
	SendText(update, ctx, out)

	// This one should be done with an uploader that stores the file id in structure, without the need to hardcode it
	//photo := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "assets/ieeehkn_logo.png")
	photo := tgbotapi.NewPhotoShare(update.Message.Chat.ID, "AgADBAADIKoxG0NL8VB8IoReowRxHedTvRkABCrHmUeeZZTrDiUBAAEC")
	ctx.Bot.Send(photo)

	out = `Questo è il nostro stemma e ogni elemento ha un significato.`
	SendText(update, ctx, out)
}
