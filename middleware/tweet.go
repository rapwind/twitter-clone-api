package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/logger"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SetTweetOnContext middleware for *gin.Context
func SetTweetOnContext(param string) gin.HandlerFunc {

	return func(c *gin.Context) {

		tid := bson.ObjectIdHex(c.Params.ByName(param))
		t, err := service.ReadTweetByID(tid)
		if err != nil {
			if err == mgo.ErrNotFound {
				errors.Send(c, errors.NotFound())
			} else {
				logger.Error(err)
				errors.Send(c, fmt.Errorf("failed to get a tweet"))
			}
			return
		}
		utils.SetTargetTweet(c, t)
	}
}
