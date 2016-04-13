package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
)

func getUser(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)

	ud, err := service.ReadUserDetailByID(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusOK, ud)
}
