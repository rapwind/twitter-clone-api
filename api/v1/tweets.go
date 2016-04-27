package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/dogenzaka/gin-tools/validation/validator"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/jsonschema"
	"github.com/techcampman/twitter-d-server/logger"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2/bson"
)

func createTweet(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors.Send(c, errors.RequestEntityTooLarge())
		return
	}
	if len(b) == 0 {
		errors.Send(c, errors.Unauthorized())
		return
	}
	if err := jsonschema.JSONSchema(b, jsonschema.V1PostTweetDocument); err != nil {
		errors.Send(c, err)
		return
	}
	t := new(entity.Tweet)
	if err = json.Unmarshal(b, t); err != nil {
		errors.Send(c, errors.BadParams("body", fmt.Sprint(b)))
		return
	}

	// get session
	uid, err := utils.GetLoginUserID(c)
	if err != nil {
		errors.Send(c, errors.Unauthorized())
		return
	}

	// Check duplicated tweets
	t0, err := service.ReadLatestTweet(uid)
	if t0 != nil && t0.Text == t.Text {
		errors.Send(c, errors.DataConflict())
		return
	}

	// Check the existence of the tweet corresponding to t.InReplyToTweetID.
	if t.InReplyToTweetID.Valid() {
		if _, err := service.ReadTweetByID(t.InReplyToTweetID); err != nil {
			errors.Send(c, errors.ErrDataNotFound)
			return
		}
	}

	if t.InRetweetToTweetID.Valid() {
		// Check the existence of the tweet corresponding to t.InRetweetToTweetID.
		if _, err := service.ReadTweetByID(t.InRetweetToTweetID); err != nil {
			errors.Send(c, errors.ErrDataNotFound)
			return
		}

		// Check duplicated retweets.
		if service.CheckDupRetweet(uid, t.InRetweetToTweetID) {
			errors.Send(c, errors.DataConflict())
			return
		}

		// Check retweeting my tweet
		t1, _ := service.ReadTweetByID(t.InRetweetToTweetID)
		if t1 != nil && t1.UserID == uid {
			errors.Send(c, errors.DataConflict())
			return
		}
	} else { // if it is not a retweet:
		if t.Text == "" { // Reject an empty tweet.
			errors.Send(c, errors.BadParams("text", ""))
			return
		}
	}

	// create tweet
	t.UserID = uid
	err = service.CreateTweet(t)
	if err != nil {
		errors.Send(c, errors.BadParams("body", fmt.Sprint(b)))
		return
	}

	if t.InReplyToTweetID.Valid() {
		// create notification
		if err := service.CreateReplyNotification(uid, t); err != nil {
			logger.Error(err)
		}
	}

	if t.InRetweetToTweetID.Valid() {
		// create notification
		if err := service.CreateRetweetNotification(uid, t); err != nil {
			logger.Error(err)
		}
	}

	c.JSON(http.StatusOK, t)
}

func deleteTweet(c *gin.Context) {
	// get params
	uid, t := getLoginUserIDAndTargetTweet(c)

	// check tweet
	t, err := service.ReadTweetByID(t.ID)
	if err != nil {
		errors.Send(c, errors.NotFound())
		return
	}
	if t.UserID != uid {
		errors.Send(c, errors.Forbidden())
		return
	}

	// delete tweet
	if err := service.RemoveTweet(t); err != nil {
		if errors.IsDataNotFound(err) {
			errors.Send(c, errors.NotFound())
		} else {
			errors.Send(c, fmt.Errorf("failed to delete a tweet"))
		}
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func getTweets(c *gin.Context) {
	// get session
	userID, _ := utils.GetLoginUserID(c)

	// Get parameters
	_, limit := utils.GetRangeParams(c, constant.DefaultLimitGetTweets)
	maxID := utils.GetObjectIDParam(c, "maxId")
	sinceID := utils.GetObjectIDParam(c, "sinceId")
	following, _ := validator.Bool{}.Check(c.Request.FormValue("following"))
	q := c.Request.FormValue("q")

	ts, err := service.ReadTweetDetails(limit, maxID, sinceID, userID, following, q)
	if err != nil {
		errors.Send(c, err)
		return
	}

	if ts == nil {
		ts = []*entity.TweetDetail{}
	}

	c.JSON(http.StatusOK, ts)
}

func getTweet(c *gin.Context) {
	loginUserID, t := getLoginUserIDAndTargetTweet(c)

	td, err := service.ReadTweetDetailByID(t.ID, loginUserID)
	if err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusOK, td)
}

func doLike(c *gin.Context) {
	loginUserID, t := getLoginUserIDAndTargetTweet(c)

	l := &entity.Like{
		UserID:  loginUserID,
		TweetID: t.ID,
	}
	if err := service.CreateLike(l); err != nil {
		errors.Send(c, err)
		return
	}

	// create notification
	if err := service.CreateLikeNotification(loginUserID, l); err != nil {
		logger.Error(err)
	}

	c.JSON(http.StatusCreated, l)
}

func undoLike(c *gin.Context) {
	loginUserID, t := getLoginUserIDAndTargetTweet(c)

	l := &entity.Like{
		UserID:  loginUserID,
		TweetID: t.ID,
	}
	if err := service.RemoveLike(l); err != nil {
		errors.Send(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func getLoginUserIDAndTargetTweet(c *gin.Context) (loginUserID bson.ObjectId, param *entity.Tweet) {
	loginUserID, _ = utils.GetLoginUserID(c)
	param, _ = utils.GetTargetTweet(c)
	return
}
