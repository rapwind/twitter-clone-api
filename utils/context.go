package utils

import (
	"fmt"

	"github.com/dogenzaka/gin-tools/validation/validator"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"gopkg.in/mgo.v2/bson"
)

// GetObjectIDPath gets ObjectID from path parameters
func GetObjectIDPath(c *gin.Context, name string) (id bson.ObjectId) {
	if v := c.Params.ByName("id"); v != "" && bson.IsObjectIdHex(v) {
		id = bson.ObjectIdHex(v)
	}
	return
}

// SetLoginUserID sets userId session on session
func SetLoginUserID(c *gin.Context, userID bson.ObjectId) {
	c.Set(constant.LoginUserKey, userID)
}

// GetRangeParams obtains "limit" and "offset" parameters
func GetRangeParams(c *gin.Context, defaultLimit int) (offset int, limit int) {
	offset, _ = validator.UInt{}.Check(c.Request.FormValue("offset"))
	limit, _ = validator.UInt{}.Check(c.Request.FormValue("limit"))
	if limit == 0 {
		limit = defaultLimit
	}
	return
}

// GetObjectIDParam gets ObjectID from query parameters
func GetObjectIDParam(c *gin.Context, name string) (id bson.ObjectId) {

	if v := c.Request.FormValue(name); v != "" && bson.IsObjectIdHex(v) {
		id = bson.ObjectIdHex(v)
	}

	return
}

// GetPoppoHeader ... get beaut header
func GetPoppoHeader(c *gin.Context) *entity.PoppoHeader {
	h := c.Request.Header
	return &entity.PoppoHeader{
		SessionID:  h.Get(constant.XPoppoSessionID),
		CSRFToken:  h.Get(constant.XPoppoCSRFToken),
		AppVersion: h.Get(constant.XPoppoAppVersion),
	}
}

// GetLoginUserID gets login userId from session
func GetLoginUserID(c *gin.Context) (userID bson.ObjectId, err error) {

	v, exists := c.Get(constant.LoginUserKey)
	if !exists {
		err = fmt.Errorf("not found in gin.Context. key = %s", constant.LoginUserKey)
		return
	}

	var ok bool
	userID, ok = v.(bson.ObjectId)
	if !ok || !userID.Valid() {
		err = fmt.Errorf("not bson.ObjectId type, value = %v", v)
	}

	return
}

// GetTargetUser ... get target user from session
func GetTargetUser(c *gin.Context) (u *entity.User, err error) {
	i, exists := c.Get(constant.TargetUserKey)
	if !exists {
		err = fmt.Errorf("not found in gin.Context. key = %s", constant.TargetUserKey)
		return
	}
	u = i.(*entity.User)
	return
}

// SetTargetUser ... set user on session
func SetTargetUser(c *gin.Context, u *entity.User) {
	c.Set(constant.TargetUserKey, u)
}

// GetTargetTweet ... get target tweet from session
func GetTargetTweet(c *gin.Context) (t *entity.Tweet, err error) {
	i, exists := c.Get(constant.TargetTweetKey)
	if !exists {
		err = fmt.Errorf("not found in gin.Context. key = %s", constant.TargetTweetKey)
		return
	}
	t = i.(*entity.Tweet)
	return
}

// SetTargetTweet ... set tweet on session
func SetTargetTweet(c *gin.Context, t *entity.Tweet) {
	c.Set(constant.TargetTweetKey, t)
}
