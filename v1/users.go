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
)

func registerUser(c *gin.Context) {

	// validate request
	installationID := c.Request.Header.Get(constant.XPoppoInstallationID)
	if installationID == "" {
		errors.Send(c, errors.Unauthorized())
		return
	}

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

	u := new(entity.User)
	u.Name = ur.Name
	u.ScreenName = ur.ScreenName
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
	id := utils.GetObjectIDPath(c, constant.IDKey)
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
