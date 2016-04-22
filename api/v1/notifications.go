package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/logger"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
)

func getNotifications(c *gin.Context) {
	loginUserID, _ := utils.GetLoginUserID(c)

	// Get parameters
	_, limit := utils.GetRangeParams(c, constant.DefaultLimitGetTweets)
	maxID := utils.GetObjectIDParam(c, "maxId")
	sinceID := utils.GetObjectIDParam(c, "sinceId")

	nds, err := service.ReadNotificationDetails(loginUserID, limit, maxID, sinceID)
	if err != nil {
		errors.Send(c, err)
		return
	}

	n := len(nds)
	if n > 0 {
		maxID, minID := nds[0].ID, nds[n-1].ID
		if maxID < minID {
			maxID, minID = minID, maxID
		}
		err = service.UpdateNotificationsUnread(loginUserID, maxID, minID)
		if err != nil {
			logger.Error(err)
		}
	}

	c.JSON(http.StatusOK, nds)
}
