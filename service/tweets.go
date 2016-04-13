package service

import (
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"gopkg.in/mgo.v2/bson"
)

// ReadTweetsCountsByUser gets entity.UserDetail.TweetsCount and LikesCount
func ReadTweetsCountsByUser(u entity.User) (tweetsCount int, likesCount int, err error) {
	tweets, err := collection.Tweets()
	if err != nil {
		return
	}
	defer tweets.Close()

	tweetsCount, err = tweets.Find(bson.M{"userId": u.UserID}).Count()
	if err != nil {
		return
	}

	likesCount = 0 // TODO: obtain likesCount!
	return
}
