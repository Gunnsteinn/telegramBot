package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gunnsteinn/telegramBot/client"
	"github.com/Gunnsteinn/telegramBot/domain"
	"math"
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
	err := json.Unmarshal(sponsorInfo, &sponsor)
	if err != nil {
		return ""
	}

	//chatText := `
	//			Buenos días <b>Facundo Ompre<a href="https://marketplace.axieinfinity.com/profile/inventory/axie">.</a></b>!!!
	//
	//			- Este es el informe de tus equipos:
	//
	//				<code>
	//				Equipo:       Geralt
	//				[%]Equipo:    100
	//				SPLs Ganados: 239
	//				</code>
	//				<code>
	//				Equipo:       Browser
	//				[%]Equipo:    100
	//				SPLs Ganados: 280
	//				</code>
	//				<code>
	//				Equipo:        Link
	//				[%]Equipo:     33
	//				SPLs Ganados:  378
	//				</code>
	//
	//			<b>Total SLP:  <i>897</i></b>
	//			<b>Total UDS:  <i>84,1386</i></b>`
	chatText := fmt.Sprintf("Buenos días <b>%s %s<a href=\"https://marketplace.axieinfinity.com/profile/inventory/axie\">.</a></b>!!!\n\n\t\t\t\t- Este es el informe de tus equipos:\n\n\t\t\t\t\t", sponsor.Name, sponsor.LastName)
	var teamSlice []string
	TotalSlp := 0
	for _, team := range sponsor.Teams {
		teamSlice = append(teamSlice, fmt.Sprintf("<code>\n\t\t\t\t\tEquipo:       %s\n\t\t\t\t\t[]Equipo:    %f\n\t\t\t\t\tSPLs Ganados: %d\n\t\t\t\t\t</code>\n\t\t\t\t\t", team.TeamName, team.PoolPercent, team.Adventurer.ProfitSlp))
		TotalSlp += TotalSlp + team.Adventurer.ProfitSlp
	}

	TotalSlp = int(math.RoundToEven(float64(TotalSlp / 2)))
	TotalUds := 0.8 * float64(TotalSlp)

	//"<b>Total SLP:  <i>897</i></b>\n\t\t\t\t" +
	//"<b>Total UDS:  <i>84,1386</i></b> %s", chatText)
	result := chatText + strings.Join(teamSlice, "") + fmt.Sprintf("<b>Total SLP:  <i>%d</i></b>\n\t\t\t\t<b>Total UDS:  <i>%f</i></b>", TotalSlp, TotalUds)
	return result
}
