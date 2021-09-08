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
	uriCryptoGuild  = "uri_crypto_guild"
	uriTelegramBot  = "uri_telegram_bot"
	uriBinancePrice = "uri_binance_price"
)

var (
	uriSponsor  = os.Getenv(uriCryptoGuild)
	uriTelegram = os.Getenv(uriTelegramBot)
	uriBinance  = os.Getenv(uriBinancePrice)
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
	sendMessage(webhookReqBody.Message.Chat.ID, textGenerator(sponsorInfo.Body))

	return nil, nil
}

// sayPolo takes a chatID and sends "polo" to them
func sendMessage(chatID int64, chatText string) error {
	// Create the request body struct

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

func textGenerator(sponsorInfo []byte) string {
	var sponsor domain.Sponsor
	json.Unmarshal(sponsorInfo, &sponsor)

	priceInfo, getAdvErr := client.ResponseClient.Get("https://api.binance.com/api/v3/ticker/price?symbol=SLPUSDT")
	if getAdvErr != nil {
		fmt.Println(getAdvErr)
	}

	fmt.Println(priceInfo)

	chatText := `
				Buenos días <b>Facundo Ompre<a href="https://marketplace.axieinfinity.com/profile/inventory/axie">.</a></b>!!!

				- Este es el informe de tus equipos:

					<code>
					Equipo:       Geralt
					[%]Equipo:    100
					SPLs Ganados: 239
					</code>
					<code>
					Equipo:       Browser
					[%]Equipo:    100
					SPLs Ganados: 280
					</code>
					<code>
					Equipo:        Link
					[%]Equipo:     33
					SPLs Ganados:  378
					</code>
				
				<b>Total SLP:  <i>897</i></b>
				<b>Total UDS:  <i>84,1386</i></b>`
	//fmt.Sprintf("Buenos días <b>Facundo Ompre<a href=\"https://marketplace.axieinfinity.com/profile/inventory/axie\">.</a></b>!!!\n\n\t\t\t\t %s" , chatText)
	//"- Este es el informe de tus equipos:\n\n\t\t\t\t\t" +
	//"<code>\n\t\t\t\t\t" +
	//"Equipo:       Geralt\n\t\t\t\t\t" +
	//"[%]Equipo:    100\n\t\t\t\t\t" +
	//"SPLs Ganados: 239\n\t\t\t\t\t" +
	//"</code>\n\t\t\t\t\t" +
	//"<code>\n\t\t\t\t\t" +
	//"Equipo:       Browser\n\t\t\t\t\t" +
	//"[%]Equipo:    100\n\t\t\t\t\t" +
	//"SPLs Ganados: 280\n\t\t\t\t\t" +
	//"</code>\n\t\t\t\t\t" +
	//"<code>\n\t\t\t\t\t" +
	//"Equipo:        Link\n\t\t\t\t\t" +
	//"[%]Equipo:     33\n\t\t\t\t\t" +
	//"SPLs Ganados:  378\n\t\t\t\t\t" +
	//"</code>\n\t\t\t\t\n\t\t\t\t" +
	//"<b>Total SLP:  <i>897</i></b>\n\t\t\t\t" +
	//"<b>Total UDS:  <i>84,1386</i></b> %s", chatText)
	return chatText
}
