package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/env"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/logger"
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

		u, err := getSession(c, sessionID)
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

// getSession gets session on a cache
func getSession(c *gin.Context, sessionID string) (userID bson.ObjectId, err error) {

	reply, err := env.GetCache().Get(constant.UserSessionPrefix + sessionID)
	if err != nil {
		if !errors.IsDataNotFound(err) {
			logger.Error(err)
		}
		return
	}
	userID = bson.ObjectIdHex(string(reply.([]byte)))

	return
}
