package service

import (
	"fmt"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CreateTweet creates "entity.Tweet" data
func CreateTweet(t *entity.Tweet) (err error) {
	if t.ID != "" {
		return fmt.Errorf("already objectId, oid = %s", t.ID)
	}

	t.ID = bson.NewObjectId()
	t.CreatedAt = t.ID.Time()

	tweets, err := collection.Tweets()
	if err != nil {
		return
	}
	defer tweets.Close()

	err = tweets.Insert(t)
	if err != nil && !mgo.IsDup(err) {
		logger.Error(err)
	}

	return
}

// RemoveTweet deletes a document on tweets collection
func RemoveTweet(t *entity.Tweet) (err error) {
	tweets, err := collection.Tweets()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer tweets.Close()

	err = tweets.RemoveId(t.ID)
	if err != nil && err != mgo.ErrNotFound {
		logger.Error(err)
	}

	return
}

// ReadTweetDetails returns an array of TweetDetail(s)
func ReadTweetDetails(limit int, maxID bson.ObjectId, userID bson.ObjectId, following bool, q string) (tds []entity.TweetDetail, err error) {
	ts, err := readTweets(limit, maxID, userID, following, q)

	tds = make([]entity.TweetDetail, len(ts))
	tdp := (*entity.TweetDetail)(nil)
	for i, t := range ts {
		tdp, err = readTweetDetailByTweet(t)
		if err != nil {
			return
		}
		tds[i] = *tdp
	}
	return
}

func readTweets(limit int, maxID bson.ObjectId, userID bson.ObjectId, following bool, q string) (ts []entity.Tweet, err error) {
	m := []bson.M{
		bson.M{"deletedAt": bson.M{"$exists": false}},
	}

	// if maxId is set:
	if maxID.Valid() {
		m = append(m, bson.M{"_id": bson.M{"$lte": maxID}})
	}

	// if query is set:
	if len(q) > 0 {
		m = append(m, bson.M{"text": bson.RegEx{Pattern: q, Options: "i"}})
	}

	// if following is set:
	if userID.Valid() && following {
		flws := []entity.Follow{}
		flws, err = ReadFollowingByID(userID, 0, -1)
		if err != nil {
			return
		}

		// Convert an array of Follow(s) into an array of user IDs
		if len(flws) != 0 {
			ids := make([]bson.ObjectId, len(flws))
			for i, v := range flws {
				ids[i] = v.TargetID
			}
			m = append(m, bson.M{"userId": bson.M{"$in": ids}})
		}
	}

	tweets, err := collection.Tweets()
	if err != nil {
		return
	}
	defer tweets.Close()

	err = tweets.Find(bson.M{"$and": m}).Sort("-_id").Limit(limit).All(&ts)
	return
}

// ReadTweetDetailByID returns TweetDetail by tweet ID
func ReadTweetDetailByID(id bson.ObjectId) (td *entity.TweetDetail, err error) {
	t, err := ReadTweetByID(id)
	if err != nil {
		return
	}

	td, err = readTweetDetailByTweet(*t)
	return
}

func readTweetDetailByTweet(t entity.Tweet) (td *entity.TweetDetail, err error) {
	tdwr, err := readTweetDetailWithoutReplyByTweet(t)
	if err != nil {
		return
	}

	inReplyToTweet := (*entity.TweetDetailWithoutReply)(nil)
	if t.InReplyToTweetID.Valid() {
		inReplyToTweet, err = readTweetDetailWithoutReplyByID(t.InReplyToTweetID)
		if err != nil {
			return
		}
	}

	td = &entity.TweetDetail{tdwr, inReplyToTweet}
	return
}

func readTweetDetailWithoutReplyByID(id bson.ObjectId) (tdwr *entity.TweetDetailWithoutReply, err error) {
	t, err := ReadTweetByID(id)
	if err != nil {
		return
	}

	tdwr, err = readTweetDetailWithoutReplyByTweet(*t)
	return
}

func readTweetDetailWithoutReplyByTweet(t entity.Tweet) (tdwr *entity.TweetDetailWithoutReply, err error) {
	u, err := ReadUserDetailByID(t.UserID)
	if err != nil {
		return
	}

	liked := false // TODO: obtain "liked"

	tdwr = &entity.TweetDetailWithoutReply{&t, u, &liked}
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

	tweetsCount, err = tweets.Find(bson.M{"userId": u.ID}).Count()
	if err != nil {
		return
	}

	likesCount = 0 // TODO: obtain likesCount!
	return
}
