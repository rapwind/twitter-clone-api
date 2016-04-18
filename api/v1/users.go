package v1

import (
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/jsonschema"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2/bson"
)

func registerUser(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors.Send(c, errors.RequestEntityTooLarge())
		return
	}
	if len(b) == 0 {
		errors.Send(c, errors.Unauthorized())
		return
	}
	if err := jsonschema.JSONSchema(b, jsonschema.V1PostUserDocument); err != nil {
		errors.Send(c, err)
		return
	}
	ur := new(entity.UserRegisterRequest)
	if err = json.Unmarshal(b, ur); err != nil {
		errors.Send(c, errors.BadParams("body", fmt.Sprint(b)))
		return
	}
	if ur.PhoneNumber == "" && ur.Email == "" {
		errors.Send(c, errors.BadParams("phoneNumber, email", "empty"))
		return
	}

	u := new(entity.User)
	u.Name = ur.Name
	u.ScreenName = ur.ScreenName
	u.Email = ur.Email
	u.PhoneNumber = utils.PhoneNumberNormalization(ur.PhoneNumber)
	u.PasswordHash = ur.PasswordHash

	// create user account
	err = service.CreateUser(u)
	if err != nil {
		errors.Send(c, errors.BadParams("screenName", u.ScreenName))
		return
	}

	c.JSON(http.StatusOK, u)
}

func getUser(c *gin.Context) {
	_, id := getLoginUserIDAndTargetUser(c)

	ud, err := service.ReadUserDetailByID(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusOK, ud)
}

func doFollow(c *gin.Context) {
	l, u := getLoginUserIDAndTargetUser(c)

	f := &entity.Follow{
		UserID:   l,
		TargetID: u.ID,
	}
	if err := service.CreateFollow(f); err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusCreated, f)
}

func undoFollow(c *gin.Context) {
	l, u := getLoginUserIDAndTargetUser(c)

	f := &entity.Follow{
		UserID:   l,
		TargetID: u.ID,
	}
	if err := service.RemoveFollow(f); err != nil {
		errors.Send(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func getFollowing(c *gin.Context) {
	_, id := getLoginUserIDAndTargetUser(c)
	offset, limit := utils.GetRangeParams(c, constant.DefaultLimitGetFollows)

	flws, err := service.ReadFollowingByID(id, offset, limit)
	if err != nil {
		errors.Send(c, err)
		return
	}

	users := make([]*entity.UserDetail, len(flws))
	for i, v := range flws {
		users[i], err = service.ReadUserDetailByID(v.TargetID)
		if err != nil {
			errors.Send(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, users)
}

func getFollower(c *gin.Context) {
	_, id := getLoginUserIDAndTargetUser(c)
	offset, limit := utils.GetRangeParams(c, constant.DefaultLimitGetFollows)

	flws, err := service.ReadFollowerByID(id, offset, limit)
	if err != nil {
		errors.Send(c, err)
		return
	}

	users := make([]*entity.UserDetail, len(flws))
	for i, v := range flws {
		users[i], err = service.ReadUserDetailByID(v.UserID)
		if err != nil {
			errors.Send(c, err)
			return
		}
	}

	c.JSON(http.StatusOK, users)
}

func getUserTweets(c *gin.Context) {
	uid := utils.GetObjectIDPath(c, constant.IDKey)

	// Get parameters
	_, limit := utils.GetRangeParams(c, constant.DefaultLimitGetTweets)
	maxID := utils.GetObjectIDParam(c, "maxId")

	ts, err := service.ReadUserTweetDetails(uid, limit, maxID)
	if err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusOK, ts)
}

func getLoginUserIDAndTargetUser(c *gin.Context) (loginUserID bson.ObjectId, param *entity.User) {
	loginUserID, _ = utils.GetLoginUserID(c)
	param, _ = utils.GetTargetUser(c)
	return
}
