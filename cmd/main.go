package main

import (
	"context"
	"log"
	"os"

	"github.com/royge/movidersms"
)

func main() {
	var (
		apikey    = os.Getenv("MOVIDER_API_KEY")
		apisecret = os.Getenv("MOVIDER_API_SECRET")
		recipient = os.Getenv("RECIPIENT")
	)

	sender := movidersms.NewSender(
		movidersms.Credentials{APIKey: apikey, APISecret: apisecret},
		[]string{},
	)
	_, err := sender.SendMessage(
		context.Background(),
		[]string{recipient},
		"TEST 1234 QWERTY",
	)
	if err != nil {
		log.Printf("unable to send message: %v", err)
	}
}
