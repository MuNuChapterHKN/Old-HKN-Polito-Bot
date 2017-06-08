package bot

import (
    "io/ioutil"
	"gopkg.in/telegram-bot-api.v4"
    "log"
	"encoding/json"
)

type JsonEvents struct {
	Events []struct {
		Uid string `json:"uid"`
		Name string `json:"name"`
		Date string `json:"date"`
		Description string `json:"description"`
        Eventbrite string `json:"eventbrite"`
		Structure []struct {
			Start string `json:"start"`
			End string `json:"end"`
			Description string `json:"description"`
			Speaker string `json:"speaker"`
		} `json:"structure"`
	} `json:"events"`
}

func RouteEventQuery(data []string, update tgbotapi.Update, ctx *BotContext) {
    if data[0] != "event" {
        log.Print("error routing")
        return
    }
	events := EventReader("events.json")
    for _, event := range events.Events {
       if event.Uid != data[1] {
            continue
       }
       if data[2] == "structure" {
           var out string
           speaker := ""
	   if talk.Speaker != "" {
		speaker = "<b>["+talk.Speaker+"]</b>\n"
	   }
           for _, talk := range event.Structure {
                out += "<b>"+talk.Start+"-"
                out += talk.End+"</b>\t\t"
                out += speaker
                out += "\t\t"+talk.Description+"\n\n"
           }
           msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, out)
           msg.ParseMode = "HTML"
           ctx.Bot.Send(msg)
       }
    }
}

func ShowEvents(update tgbotapi.Update, ctx *BotContext) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I prossimi eventi sono:")
	ctx.Bot.Send(msg)
	// open the events file
	events := EventReader("events.json")
    for _, event := range events.Events {
        out := "<b>" + event.Name + "</b>: Date "+ event.Date+"\n"
        out += event.Description
        msg = tgbotapi.NewMessage(update.Message.Chat.ID, out)
        button1 := tgbotapi.NewInlineKeyboardButtonData("Programma \xF0\x9F\x8E\xAB", "event:"+event.Uid+":structure")
        button2 := tgbotapi.NewInlineKeyboardButtonURL("Eventbrite \xE2\x9D\x93", event.Eventbrite)
        row := tgbotapi.NewInlineKeyboardRow(button1, button2)
        keyboard := tgbotapi.NewInlineKeyboardMarkup(row)
        msg.ReplyMarkup = keyboard
        msg.ParseMode = "HTML"
        ctx.Bot.Send(msg)
	}
}

func EventReader(fileName string) JsonEvents {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var events JsonEvents
	err = json.Unmarshal(buf, &events)
	if err != nil {
		log.Fatal(err)
	}
	return events
}
