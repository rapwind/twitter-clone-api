package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2/bson"
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
	getFollows(c, "userId",
		func(f entity.Follow) bson.ObjectId {
			return f.TargetID
		},
		constant.DefaultLimitFollowingUsers)
}

func getFollower(c *gin.Context) {
	getFollows(c, "targetId",
		func(f entity.Follow) bson.ObjectId {
			return f.UserID
		},
		constant.DefaultLimitFollowerUsers)
}

func getFollows(c *gin.Context, key string, getKey func(entity.Follow) bson.ObjectId, defaultLimit int) {
	id := utils.GetObjectIDPath(c, constant.IDKey)
	offset, limit := utils.GetRangeParams(c, defaultLimit)

	flws, err := service.ReadFollowsByID(id, key, offset, limit)
	if err != nil {
		errors.Send(c, err)
		return
	}

	users := make([]*entity.UserDetail, len(flws))
	for i, v := range flws {
		users[i], err = service.ReadUserDetailByID(getKey(v))
		if err != nil {
			errors.Send(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, users)
}
