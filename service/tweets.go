package service

import (
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"gopkg.in/mgo.v2/bson"
)

// ReadTweetDetailByID returns TweetDetail by tweet ID
func ReadTweetDetailByID(id bson.ObjectId) (td *entity.TweetDetail, err error) {
	t, err := ReadTweetDetailWithoutReplyByID(id)
	if err != nil {
		return
	}

	inReplyToTweet := (*entity.TweetDetailWithoutReply)(nil)
	if t.InReplyToTweetID.Valid() {
		inReplyToTweet, err = ReadTweetDetailWithoutReplyByID(t.InReplyToTweetID)
		if err != nil {
			return
		}
	}

	td = &entity.TweetDetail{t, inReplyToTweet}
	return
}

// ReadTweetDetailWithoutReplyByID returns entity.TweetDetailWithoutReply by tweet ID
func ReadTweetDetailWithoutReplyByID(id bson.ObjectId) (tdwr *entity.TweetDetailWithoutReply, err error) {
	t, err := ReadTweetByID(id)
	if err != nil {
		return
	}

	u, err := ReadUserDetailByID(t.UserID)
	if err != nil {
		return
	}

	liked := false // TODO: obtain "liked"

	tdwr = &entity.TweetDetailWithoutReply{t, u, &liked}
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
