package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/env"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/logger"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2/bson"
)

// CheckSession is a check session by SessionID
func CheckSession() gin.HandlerFunc {

	return func(c *gin.Context) {
		sessionID := c.Request.Header.Get(constant.XPoppoSessionID)
		if sessionID == "" {
			errors.Send(c, errors.Unauthorized())
			return
		}

		u, err := getSession(c)
		if err != nil {
			if errors.IsDataNotFound(err) {
				errors.Send(c, errors.Unauthorized())
				return
			}
			errors.Send(c, fmt.Errorf("failed to get a session"))
			return
		}

		if !u.Valid() {
			errors.Send(c, errors.Unauthorized())
			return
		}

		// set login user id for gin.Context
		utils.SetLoginUserID(c, u)
	}
}

// SetSession sets userID on session
func SetSession(c *gin.Context, userID bson.ObjectId) (err error) {

	if !userID.Valid() {
		return fmt.Errorf("invalid userID = %v", userID)
	}

	sid := string(bson.NewObjectId().Hex())
	err = env.GetCache().Set(constant.UserSessionPrefix+sid, userID.Hex(), constant.SessionExpires)
	if err != nil {
		logger.Error(err)
		return
	}

	c.Writer.Header().Set(constant.XPoppoSessionID, sid)
	return
}

// DelSession deletes session on a cache
func DelSession(c *gin.Context) (err error) {
	sid := utils.GetPoppoHeader(c).SessionID

	err = service.RemoveSessionByID(bson.ObjectIdHex(sid))
	if err != nil {
		logger.Error(err)
	}
	_, err = env.GetCache().Delete(constant.UserSessionPrefix + utils.GetPoppoHeader(c).SessionID)
	if err != nil && !errors.IsDataNotFound(err) {
		logger.Error(err)
	}

	return
}

// SetLoginUserIDIfNotEmpty sets login user id if not empty
func SetLoginUserIDIfNotEmpty() gin.HandlerFunc {

	return func(c *gin.Context) {
		ph := utils.GetPoppoHeader(c)
		if ph.SessionID == "" {
			return
		}

		u, err := getSession(c)
		if err != nil {
			return
		}

		if u.Valid() {
			// set login user id for gin.Context
			utils.SetLoginUserID(c, u)
		}
	}
}

// getSession gets session on a cache
func getSession(c *gin.Context) (userID bson.ObjectId, err error) {
	sid := utils.GetPoppoHeader(c).SessionID

	reply, err := env.GetCache().Get(constant.UserSessionPrefix + sid)
	if err != nil {
		s, err := service.ReadSessionByID(bson.ObjectIdHex(sid))
		if err == nil {
			userID = s.UserID
			SetSession(c, userID)
		}
	} else {
		userID = bson.ObjectIdHex(string(reply.([]byte)))
	}

	return
}
