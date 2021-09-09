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
	"strconv"
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
	errSponsor := json.Unmarshal(result.Body, &sponsor)
	if errSponsor != nil {
		return nil, errSponsor
	}

	return &sponsor, nil
}

func TelegramProcessorService(webhookReqBody domain.WebhookReqBody) (*domain.SendMessageReqBody, error) {
	fmt.Println("reply sent " + string(uriSponsor+strings.ToLower(webhookReqBody.Message.Text)))

	if !(len(webhookReqBody.Message.Text) > 0) {
		return nil, sendMessage(webhookReqBody.Message.Chat.ID, "Empty Text.")
	}

	sponsorInfo, getAdvErr := client.ResponseClient.Get(uriSponsor + strings.ToLower(webhookReqBody.Message.Text))
	fmt.Println("------- 1 -------")
	fmt.Println(getAdvErr)
	fmt.Println("------- 2 -------")
	fmt.Println(sponsorInfo)
	fmt.Println("------- 3 -------")
	if getAdvErr != nil {
		fmt.Println("------- 4 -------")
		return nil, sendMessage(webhookReqBody.Message.Chat.ID, "Wrong user.")
	}
	// log a confirmation message if the message is sent successfully
	//fmt.Println("reply sent" + string(sponsorInfo.Body))
	err := sendMessage(webhookReqBody.Message.Chat.ID, textGenerator(sponsorInfo.Body))
	fmt.Println("------- 5 -------")
	if err != nil {
		fmt.Println("------- 6 -------")
		return nil, err
	}
	fmt.Println("------- 7 -------")
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

	// Send a post request with your token
	res, err := http.Post(uriTelegram, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

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

	priceInfo, getAdvErr := client.ResponseClient.Get(uriBinance)
	if getAdvErr != nil {
		fmt.Println(getAdvErr)
	}

	var binancePrice domain.BinancePrice
	errBinancePrice := json.Unmarshal(priceInfo.Body, &binancePrice)
	if errBinancePrice != nil {
		fmt.Println(errBinancePrice)
	}

	chatText := fmt.Sprintf("Buenos días <b>%s %s<a href=\"https://storage.googleapis.com/assets.axieinfinity.com/axies/5684/axie/axie-full-transparent.png\">.</a></b>!!!\n\n\t\t\t\t- Este es el informe de tus equipos:\n\n\t\t\t\t\t", sponsor.Name, sponsor.LastName)
	var teamSlice []string
	TotalSlp := 0
	for _, team := range sponsor.Teams {
		sponsorProfitSlp := int(math.RoundToEven(float64(team.Adventurer.ProfitSlp / 2)))
		teamSlice = append(teamSlice, fmt.Sprintf("<code>\n\t\t\t\t\tEquipo:       %s\n\t\t\t\t\t[%s]Equipo:    %f\n\t\t\t\t\tSPLs Ganados: %d\n\t\t\t\t\t</code>\n\t\t\t\t\t", team.TeamName, "%", team.PoolPercent, sponsorProfitSlp))
		TotalSlp += TotalSlp + (sponsorProfitSlp * (int(team.PoolPercent) / 100))
	}
	price, _ := strconv.ParseFloat(binancePrice.Price, 64)
	TotalUds := price * float64(TotalSlp)

	result := chatText + strings.Join(teamSlice, "") + fmt.Sprintf("<b>Total SLP:  <i>%d</i></b>\n\t\t\t\t<b>Total UDS:  <i>%f</i></b>", TotalSlp, TotalUds)
	return result
}
