package service

import (
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"gopkg.in/mgo.v2/bson"
)

// ReadTweetDetailWithoutReplyByID returns entity.TweetDetailWithoutReply by tweet ID
func ReadTweetDetailWithoutReplyByID(id bson.ObjectId) (tdwr *entity.TweetDetailWithoutReply, err error) {
	t, err := ReadTweetByID(id)
	if err != nil {
		return
	}

	u, err := ReadUserByID(t.UserID)
	if err != nil {
		return
	}

	tdwr = &entity.TweetDetailWithoutReply{t, u}
	return
}

// ReadTweetByID returns entity.Tweet by tweet ID
func ReadTweetByID(id bson.ObjectId) (t *entity.Tweet, err error) {
	tweets, err := collection.Tweets()
	if err != nil {
		return
	}
	defer tweets.Close()

	t = new(entity.Tweet)
	err = tweets.Find(bson.M{"_id": id}).One(t)
	return
}

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
