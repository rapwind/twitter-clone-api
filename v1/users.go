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
	getFollows(c, constant.FollowingIDKey, constant.GetFollowerIDKey, constant.DefaultLimitGetFollowing)
}

func getFollower(c *gin.Context) {
	getFollows(c, constant.FollowerIDKey, constant.GetFollowingIDKey, constant.DefaultLimitGetFollower)
}

// getFollows performs as GET FOLLOWING or GET FOLLOWER API
//   originKey     ... the key of user ID: you obtain following users or followers of this user
//                     (constant.FollowingIDKey or constant.FollowerIDKey).
//   getPartnerKey ... a function that obtains the key of user IDs of following users or followers
//                     (constant.GetFollowingIDKey or constant.GetFollowerIDKey).
func getFollows(c *gin.Context, originKey constant.FollowsKey, getPartnerKey func(entity.Follow) bson.ObjectId, defaultLimit int) {
	id := utils.GetObjectIDPath(c, constant.IDKey)
	offset, limit := utils.GetRangeParams(c, defaultLimit)

	flws, err := service.ReadFollowsByQuery(bson.M{(string)(originKey): id}, offset, limit)
	if err != nil {
		errors.Send(c, err)
		return
	}

	users := make([]*entity.UserDetail, len(flws))
	for i, v := range flws {
		users[i], err = service.ReadUserDetailByID(getPartnerKey(v))
		if err != nil {
			errors.Send(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, users)
}
