package service

import (
	"github.com/techcampman/twitter-d-server/db/collection"
	"gopkg.in/mgo.v2/bson"
)

// ReadTweetsCountsByUserID gets entity.UserDetail.TweetsCount and LikesCount
func ReadTweetsCountsByUserID(id bson.ObjectId) (tweetsCount int, likesCount int, err error) {
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
