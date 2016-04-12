package utils

import (
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
