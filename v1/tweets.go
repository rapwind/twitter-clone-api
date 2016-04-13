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

func getTweet(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)

	t, err := service.ReadTweetDetailWithoutReplyByID(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	inReplyToTweet := (*entity.TweetDetailWithoutReply)(nil)
	if t.InReplyToTweetID.Valid() {
		inReplyToTweet, err = service.ReadTweetDetailWithoutReplyByID(t.InReplyToTweetID)
		if err != nil {
			errors.Send(c, err)
			return
		}
	}

	td := &entity.TweetDetail{t, inReplyToTweet}

	c.JSON(http.StatusOK, td)
}
