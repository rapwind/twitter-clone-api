package v1

import (
	"net/http"

	"github.com/dogenzaka/gin-tools/validation/validator"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
)

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
