package app

import (
	"github.com/Gunnsteinn/telegramBot/controller/ping"
	"github.com/Gunnsteinn/telegramBot/controller/telegramProcessor"
)

func mapUrls() {
	router.GET("/healthCheck", ping.HelthCheckStatus)

	router.POST("/telegramListener", telegramProcessor.TelegramProcessor)
	router.GET("/:sponsors_id", telegramProcessor.Test)
}
