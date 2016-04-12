package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2/bson"
)

func getUser(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)

	u, err := getUserInfo(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	tweetsCount, likesCount, err := getTweetsInfo(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	followingCount, followerCount, err := getFollowsInfo(id)
	if err != nil {
		errors.Send(c, err)
		return
	}

	ud := &entity.UserDetail{u, tweetsCount, likesCount, followerCount, followingCount, nil}

	c.JSON(http.StatusOK, ud)
}

func getUserInfo(id bson.ObjectId) (u *entity.User, err error) {
	users, err := collection.Users()
	if err != nil {
		return
	}
	defer users.Close()

	u = new(entity.User)
	err = users.Find(bson.M{"_id": id}).One(u)
	return
}

func getTweetsInfo(id bson.ObjectId) (tweetsCount int, likesCount int, err error) {
	tweets, err := collection.Tweets()
	if err != nil {
		return
	}
	defer tweets.Close()

	tweetsCount, err = tweets.Find(bson.M{"userId": id}).Count()
	if err != nil {
		return
	}

	likesCount = 0 // TODO: obtain likesCount!
	return
}

func getFollowsInfo(id bson.ObjectId) (followingCount int, followerCount int, err error) {
	follows, err := collection.Follows()
	if err != nil {
		return
	}
	defer follows.Close()

	followingCount, err = follows.Find(bson.M{"userId": id}).Count()
	if err != nil {
		return
	}

	followerCount, err = follows.Find(bson.M{"targetId": id}).Count()
	return
}
