package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
	"fmt"
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

func getFollowing(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)
	offset, limit := utils.GetRangeParams(c, constant.DefaultLimitFollowingUsers)

	flws, err := service.ReadFollowingByID(id, offset, limit)
	if err != nil {
		errors.Send(c, err)
		return
	}

	fmt.Println(flws)

	users := make([]*entity.UserDetail, limit)
	for i, v := range flws {
		users[i], err = service.ReadUserDetailByID(v.TargetID)
		if err != nil {
			fmt.Println(v.TargetID)
			errors.Send(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, users)
}
