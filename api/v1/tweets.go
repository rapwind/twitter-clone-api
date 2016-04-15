package v1

import (
	"net/http"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dogenzaka/gin-tools/validation/validator"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/jsonschema"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
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
		errors.Send(c, fmt.Errorf("failed to get a login user"))
		return
	}

	// create tweet
	t.UserID = uid
	err = service.CreateTweet(t)
	if err != nil {
		errors.Send(c, errors.BadParams("body", fmt.Sprint(b)))
		return
	}

	c.JSON(http.StatusOK, t)
}

func getTweets(c *gin.Context) {
	_, limit := utils.GetRangeParams(c, constant.DefaultLimitGetTweets)
	maxID := utils.GetObjectIDParam(c, "maxId")
	userID := utils.GetObjectIDParam(c, "userId")
	following, _ := validator.Bool{}.Check(c.Request.FormValue("following"))
	q := c.Request.FormValue("q")

	ts, err := service.ReadTweetDetails(limit, maxID, userID, following, q)
	if err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusOK, ts)
}

func getTweet(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)

	td, err := service.ReadTweetDetailByID(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	c.JSON(http.StatusOK, td)
}
