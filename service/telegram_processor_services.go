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
	chatText = `- Negrita: <b>texto en negrita</b>, <strong>negrita</strong>
				- Cursiva: <i>texto en cursiva</i>, <em>cursiva</em>

				- Subrayado: <u>texto subrayado</u>, <ins>subrayado</ins>

				- Tachado: <s>texto tachado</s>, <strike>tachado</strike>, <del>tachado</del>

				Enlaces: <a href="https://www.Tecnonucleous.com/">Tecnonucleous</a>

				<a href="https://www.carspecs.us/photos/c8447c97e355f462368178b3518367824a757327-2000.jpg"> </a>

				<a href="https://storage.googleapis.com/assets.axieinfinity.com/axies/3624156/axie/axie-full-transparent.png">&#8205;</a>

				<a href="tg://user?id=123456789">inline mention of a user</a>
				- Código: <code>texto con el código</code>

				- Indica a la API que debe respetar los saltos de línea y los espacios en blanco: <pre>Texto con saltos de línea y espacios</pre>`

	reqBody := domain.SendMessageReqBody{
		ChatID:    chatID,
		Text:      chatText,
		ParseMode: "HTML",
	}

	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	fmt.Println(reqBody)
	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot1913861473:AAGT0ranx9RBMrtRVzrLx5PYiakOsNH6VOE/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	fmt.Println(res)
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
