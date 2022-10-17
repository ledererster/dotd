package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("you must include a c g or t")
	}
	store := os.Args[1]
	var title, data string
	var tries int
	for {
		switch store {
		case "c":
			title, data = getCardHausDeal()
		case "g":
			title, data = getGameNerdzDeal()
		case "t":
			title, data = getTabletopMerchantDotd()

		}
		title = strings.TrimSpace(title)
		if diff(store, title) {
			break
		}
		if tries += 1; tries > 5 {
			log.Fatalf("too many tries. not updated")
		}
		time.Sleep(1 * time.Minute)
	}
	bggLink := searchBGG(title)
	token := getToken()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panicf("failed to load bot : %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(504504957, data+"\n\n"+bggLink)

	_, err = bot.Send(msg)
	if err != nil {
		log.Panicf("failed to send message: %v", err)
	}

}

func getToken() string {
	dat, err := os.ReadFile("telegram.token")
	if err != nil {
		log.Panicf("can't read token!: %v", err)
	}

	return string(dat)
}
