package utils

import (
	"github.com/dogenzaka/gin-tools/validation/validator"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
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
