package telegramProcessor

import (
	"fmt"
	"github.com/Gunnsteinn/telegramBot/domain"
	"github.com/Gunnsteinn/telegramBot/service"
	"github.com/Gunnsteinn/telegramBot/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TelegramProcessor(c *gin.Context) {
	var webhookReqBody domain.WebhookReqBody
	if err := c.ShouldBindJSON(&webhookReqBody); err != nil {
		restErr := errors.NewBadRequestError("invalid json body.")
		c.JSON(restErr.Status, restErr)
		return
	}
	fmt.Println(webhookReqBody)
	result, getAdvErr := service.TelegramProcessorService(webhookReqBody)
	if getAdvErr != nil {
		c.JSON(http.StatusInternalServerError, getAdvErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func Test(c *gin.Context) {
	result, getAdvErr := service.GetSponsor(c.Param("sponsors_id"))
	if getAdvErr != nil {
		c.JSON(http.StatusInternalServerError, getAdvErr)
		return
	}

	c.JSON(http.StatusOK, result)
}
