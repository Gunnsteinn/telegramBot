package telegramProcessor

import (
	"encoding/hex"
	"fmt"
	"github.com/Gunnsteinn/telegramBot/domain"
	"github.com/Gunnsteinn/telegramBot/service"
	"github.com/Gunnsteinn/telegramBot/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
)

func getSponsorId(sponsorWalletParam string) (string, *errors.RestErr) {
	if _, addressHexErr := hex.DecodeString(sponsorWalletParam); addressHexErr == nil {
		return sponsorWalletParam, nil
	}
	if _, addressMailErr := mail.ParseAddress(sponsorWalletParam); addressMailErr == nil {
		return sponsorWalletParam, nil
	}
	return "", errors.NewBadRequestError("user id should be a hex or mail.")
}

func TelegramProcessor(c *gin.Context) {
	fmt.Println("Controller 1")
	var webhookReqBody domain.WebhookReqBody
	if err := c.ShouldBindJSON(&webhookReqBody); err != nil {
		restErr := errors.NewBadRequestError("invalid json body.")
		c.JSON(restErr.Status, restErr)
		return
	}
	fmt.Println("Controller 2")
	result, getAdvErr := service.TelegramProcessorService(webhookReqBody)
	if getAdvErr != nil {
		c.JSON(http.StatusInternalServerError, getAdvErr)
		return
	}
	fmt.Println("Controller 3")
	c.JSON(http.StatusOK, result)
}

func Test(c *gin.Context) {
	sponsorID, idErr := getSponsorId(c.Param("sponsors_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	result, getAdvErr := service.GetSponsor(sponsorID)
	if getAdvErr != nil {
		c.JSON(http.StatusInternalServerError, getAdvErr)
		return
	}

	c.JSON(http.StatusOK, result)
}
