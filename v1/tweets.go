package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
)

func getTweet(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)

	td, err := service.ReadTweetDetailByID(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusOK, td)
}
