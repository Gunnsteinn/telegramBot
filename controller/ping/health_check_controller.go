package ping

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelthCheckStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Status": "UP",
	})
}
