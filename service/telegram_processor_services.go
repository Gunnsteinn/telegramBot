package service

import (
	"encoding/json"
	"fmt"
	"github.com/Gunnsteinn/telegramBot/client"
	"github.com/Gunnsteinn/telegramBot/domain"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	uriCryptoGuild  = "uri_crypto_guild"
	uriTelegramBot  = "uri_telegram_bot"
	uriBinancePrice = "uri_binance_price"
	sunEmoji        = "\xF0\x9F\x8C\x9E"
	moonEmoji       = "\xF0\x9F\x8C\x9C"
)

var (
	uriSponsor  = os.Getenv(uriCryptoGuild)
	uriTelegram = os.Getenv(uriTelegramBot)
	uriBinance  = os.Getenv(uriBinancePrice)
	axiesArray  = []string{"2679", "5684", "2579", "2183", "5261", "1578", "2336", "1337", "1889", "1301", "2403"}
)

func GetSponsor(sponsorId string) (*domain.Sponsor, error) {
	result, err := client.ResponseClient.Get(uriSponsor + "?filter=" + sponsorId)
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
	if !(len(webhookReqBody.Message.Text) > 0) {
		return nil, sendMessage(webhookReqBody.Message.Chat.ID, "Empty Text.")
	}

	sponsorNickName := sponsorValidation(strings.ToLower(webhookReqBody.Message.Text))
	if len(sponsorNickName) < 1 {
		return nil, sendMessage(webhookReqBody.Message.Chat.ID, "Wrong user.")
	}

	sponsorInfo, getAdvErr := client.ResponseClient.Get(uriSponsor + "?filter=" + sponsorNickName)
	if getAdvErr != nil || sponsorInfo.StatusCode >= 400 {
		return nil, sendMessage(webhookReqBody.Message.Chat.ID, "Wrong user.")
	}

	err := sendMessage(webhookReqBody.Message.Chat.ID, textGenerator(sponsorInfo.Body))
	if err != nil {
		return nil, sendMessage(webhookReqBody.Message.Chat.ID, "Wrong sendMessage.")
	}

	return nil, nil
}

func sendMessage(chatID int64, chatText string) error {
	result, getAdvErr := client.ResponseClient.Post(uriTelegram, domain.SendMessageReqBody{
		ChatID:    chatID,
		Text:      chatText,
		ParseMode: "HTML",
	})
	if getAdvErr != nil {
		fmt.Println(getAdvErr)
	}

	fmt.Println(result.Status)
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
		binancePrice.Price = "0"
	}

	return formatChatText(sponsor, binancePrice)
}

func sponsorValidation(webhookReqBodyMessageText string) string {
	Aux := strings.Split(webhookReqBodyMessageText, "/")

	if len(Aux) < 2 {
		return ""
	}

	if len(Aux[1]) < 1 {
		return ""
	}

	return Aux[1]
}

func formatChatText(sponsor domain.Sponsor, binancePrice domain.BinancePrice) string {
	n := rand.Int() % len(axiesArray)
	emoticon := sunEmoji
	if time.Now().Hour() > 23 || time.Now().Hour() < 11 {
		emoticon = moonEmoji
	}

	chatText := fmt.Sprintf("Hola <b>%s %s<a href=\"https://storage.googleapis.com/assets.axieinfinity.com/axies/%s/axie/axie-full-transparent.png\">.</a></b>!!! \t\t\t...%s \n\n\t\t\t\t- Este es el informe de tus equipos:\n\n\t\t\t\t\t", sponsor.Name, sponsor.LastName, axiesArray[n], emoticon)
	var teamSlice []string
	TotalSlp := 0
	for _, team := range sponsor.Teams {
		sponsorProfitSlp := int(math.Round(float64(team.Adventurer.ProfitSlp)/2) * (team.PoolPercent / 100))
		teamSlice = append(teamSlice, fmt.Sprintf("<code>\n\t\t\t\tEquipo:       %s\n\t\t\t\t[%s]Equipo:    %d\n\t\t\t\tSLPs Ganados: %d\n\t\t\t\t\t</code>\n\t\t\t\t\t", team.TeamName, "%", int(team.PoolPercent), sponsorProfitSlp))
		TotalSlp += sponsorProfitSlp
	}

	price, _ := strconv.ParseFloat(binancePrice.Price, 64)
	TotalUds := price * float64(TotalSlp)

	return chatText + strings.Join(teamSlice, "") + fmt.Sprintf("<b>Total SLP:   <i>%d</i></b>\n\t\t\t\t<b> SLP/USDT:  <i>%f</i></b>\n\t\t\t\t\t<b>Total UDS:   <i>%f</i></b>\n\t\t\t\t\t", TotalSlp, price, TotalUds)
}
