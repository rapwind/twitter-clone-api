package v1

import (
	"fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2/bson"
)

func iso8601(s string) (t time.Time) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		fmt.Println("Oops!")
	}
	return
}

func getUser(c *gin.Context) {
	id := utils.GetObjectIDPath(c, constant.IDKey)

	users, err := collection.Users()
	if err != nil {
		errors.Send(c, err)
		return
	}
	defer users.Close()

	// Test Data
	/*users.Insert(&entity.User{
		UserID: bson.NewObjectId(), // 570c7bcd302125232a3021ca
		Name: "ほげふが",
		ScreenName: "hogehuga",
		ProfileImageURL: "https://pbs.twimg.com/profile_images/666407537084796928/YBGgi9BO.png",
		ProfileBackgroundImageURL: "https://pbs.twimg.com/profile_banners/2931934340/1456564386",
		Biography: "よろしくおねがいします。",
		LocationText: "日本, 渋谷, 道玄坂",
		URL: "https://www.cyberagent.co.jp/",
		Birthday: nil,
		CreatedAt: iso8601("2016-03-20T12:34:56+09:00"),
		UpdatedAt: iso8601("2016-04-11T12:34:56+09:00"),
	})*/

	u := new(entity.User)
	err = users.Find(bson.M{"_id": id}).One(u)
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
	if err != nil {
		return
	}
	return
}
