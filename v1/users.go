package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
)

func getUser(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)

	u, err := service.GetUser(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	tweetsCount, likesCount, err := service.GetTweetsInfo(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	followingCount, followerCount, err := service.GetFollowsInfo(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	ud := &entity.UserDetail{u, tweetsCount, likesCount, followerCount, followingCount, nil}

	c.JSON(http.StatusOK, ud)
}
