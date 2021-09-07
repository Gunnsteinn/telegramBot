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

	sponsorInfo, getAdvErr := client.ResponseClient.Get(uriSponsor + strings.ToLower(webhookReqBody.Message.Text))
	if getAdvErr != nil {
		sendMessage(webhookReqBody.Message.Chat.ID, "Manco")
	}
	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent" + string(sponsorInfo.Body))
	sendMessage(webhookReqBody.Message.Chat.ID, "")

	return nil, nil
}

// sayPolo takes a chatID and sends "polo" to them
func sendMessage(chatID int64, chatText string) error {
	// Create the request body struct
	chatText = "*bold *text* _italic *text_ __underline__ " +
		"~strikethrough~ *bold _italic bold " +
		"~italic bold strikethrough~ __underline italic bold___ bold* [inline URL](http://www.example.com/) " +
		"[inline mention of a user](tg://user?id=123456789) " +
		"`inline fixed-width code` " +
		"```pre-formatted fixed-width code block``` ```python " +
		"pre-formatted fixed-width code block written in the Python programming language" +
		"```"
	reqBody := domain.SendMessageReqBody{
		ChatID: chatID,
		Text:   chatText,
	}

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot1913861473:AAGT0ranx9RBMrtRVzrLx5PYiakOsNH6VOE/sendMessage?parse_mode=MarkdownV2", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
