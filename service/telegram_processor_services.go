package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gunnsteinn/telegramBot/client"
	"github.com/Gunnsteinn/telegramBot/domain"
	"net/http"
	"os"
	"strings"
)

const (
	uriCryptoGuild = "uri_crypto_guild"
	uriTelegramBot = "uri_telegram_bot"
)

var (
	uriSponsor  = os.Getenv(uriCryptoGuild)
	uriTelegram = os.Getenv(uriTelegramBot)
)

func GetSponsor(sponsorId string) (*domain.Sponsor, error) {
	result, err := client.ResponseClient.Get(uriSponsor + sponsorId)
	if err != nil {
		return nil, err
	}

	var sponsor domain.Sponsor
	json.Unmarshal(result.Body, &sponsor)

	return &sponsor, nil
}

func TelegramProcessorService(webhookReqBody domain.WebhookReqBody) (*domain.SendMessageReqBody, error) {

	//sponsorInfo, getAdvErr := GetSponsor())
	//if getAdvErr != nil {
	//	sendMessage(webhookReqBody.Message.Chat.ID, "Manco")
	//}
	fmt.Println("Service 1")
	sponsorInfo, getAdvErr := client.ResponseClient.Get(uriSponsor + strings.ToLower(webhookReqBody.Message.Text))
	if getAdvErr != nil {
		sendMessage(webhookReqBody.Message.Chat.ID, "Manco")
	}
	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
	sendMessage(webhookReqBody.Message.Chat.ID, string(sponsorInfo.Body))

	return nil, nil
}

// sayPolo takes a chatID and sends "polo" to them
func sendMessage(chatID int64, chatText string) error {
	// Create the request body struct
	reqBody := domain.SendMessageReqBody{
		ChatID: chatID,
		Text:   chatText,
	}
	fmt.Println("Service 2")
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	fmt.Println("Service 3")
	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot1913861473:AAGT0ranx9RBMrtRVzrLx5PYiakOsNH6VOE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}