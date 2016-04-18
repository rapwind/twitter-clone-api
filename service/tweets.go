package service

import (
	"fmt"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
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

// ReadUserTweetDetails returns an array of TweetDetail(s) posted by a given user.
func ReadUserTweetDetails(userID bson.ObjectId, loginUserID bson.ObjectId, limit int, maxID bson.ObjectId) (tds []entity.TweetDetail, err error) {
	m := []bson.M{
		bson.M{"userId": userID},
		bson.M{"deletedAt": bson.M{"$exists": false}},
	}

	// if maxId is set:
	if maxID.Valid() {
		m = append(m, bson.M{"_id": bson.M{"$lte": maxID}})
	}

	tds, err = readSortedTweetDetails(bson.M{"$and": m}, limit, loginUserID)
	return
}

// ReadTweetDetails returns an array of TweetDetail(s) (userID = login user ID)
func ReadTweetDetails(limit int, maxID bson.ObjectId, userID bson.ObjectId, following bool, q string) (tds []entity.TweetDetail, err error) {
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

	tds, err = readSortedTweetDetails(bson.M{"$and": m}, limit, userID)
	return
}

// ReadTweetDetailByID returns TweetDetail by tweet ID
func ReadTweetDetailByID(id bson.ObjectId, loginUserID bson.ObjectId) (td *entity.TweetDetail, err error) {
	t, err := ReadTweetByID(id)
	if err != nil {
		return
	}

	td, err = readTweetDetailByTweet(*t, loginUserID)
	return
}

func readSortedTweetDetails(m bson.M, limit int, loginUserID bson.ObjectId) (tds []entity.TweetDetail, err error) {
	tweets, err := collection.Tweets()
	if err != nil {
		return
	}
	defer tweets.Close()

	ts := []entity.Tweet{}
	err = tweets.Find(m).Sort("-createdAt").Limit(limit).All(&ts)
	if err != nil {
		return
	}

	tds, err = readTweetsDetailByTweets(ts, loginUserID)
	return
}

func readTweetsDetailByTweets(ts []entity.Tweet, loginUserID bson.ObjectId) (tds []entity.TweetDetail, err error) {
	tds = make([]entity.TweetDetail, len(ts))
	tdp := (*entity.TweetDetail)(nil)
	for i, t := range ts {
		tdp, err = readTweetDetailByTweet(t, loginUserID)
		if err != nil {
			return
		}
		tds[i] = *tdp
	}
	return
}

func readTweetDetailByTweet(t entity.Tweet, loginUserID bson.ObjectId) (td *entity.TweetDetail, err error) {
	tdwr, err := readTweetDetailWithoutReplyByTweet(t, loginUserID)
	if err != nil {
		return
	}

	inReplyToTweet := (*entity.TweetDetailWithoutReply)(nil)
	if t.InReplyToTweetID.Valid() {
		inReplyToTweet, err = readTweetDetailWithoutReplyByID(t.InReplyToTweetID, loginUserID)
		if err != nil {
			return
		}
	}

	td = &entity.TweetDetail{tdwr, inReplyToTweet}
	return
}

func readTweetDetailWithoutReplyByID(tweetID bson.ObjectId, loginUserID bson.ObjectId) (tdwr *entity.TweetDetailWithoutReply, err error) {
	t, err := ReadTweetByID(tweetID)
	if err != nil {
		return
	}

	tdwr, err = readTweetDetailWithoutReplyByTweet(*t, loginUserID)
	return
}

func readTweetDetailWithoutReplyByTweet(t entity.Tweet, loginUserID bson.ObjectId) (tdwr *entity.TweetDetailWithoutReply, err error) {
	u, err := ReadUserDetailByID(t.UserID)
	if err != nil {
		return
	}

	likes, err := collection.Likes()
	if err != nil {
		return
	}
	defer likes.Close()

	var liked bool
	l := entity.Like{}
	if loginUserID.Valid() && likes.Find(bson.M{"userId": loginUserID, "tweetId": t.ID}).One(&l) == nil {
		liked = true
	} else {
		liked = false
	}

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

	likes, err := collection.Likes()
	if err != nil {
		return
	}
	defer likes.Close()

	likesCount, err = likes.Find(bson.M{"userId": u.ID}).Count()
	return
}

// CreateLike creates "entity.Like" data
func CreateLike(l *entity.Like) (err error) {
	if err := checkValidLike(l); err != nil {
		return errors.BadParams("like", "invalid")
	}

	if checkLiked(l) {
		err = errors.DataConflict()
		return
	}

	l.ID = bson.NewObjectId()
	l.CreatedAt = l.ID.Time()

	likes, err := collection.Likes()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer likes.Close()

	err = likes.Insert(l)
	if err != nil {
		logger.Error(err)
	}

	return
}

// RemoveLike deletes a document on a collection
func RemoveLike(l *entity.Like) (err error) {
	if err := checkValidLike(l); err != nil {
		return errors.BadParams("like", "invalid")
	}

	likes, err := collection.Likes()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer likes.Close()

	err = likes.Remove(bson.M{"userId": l.UserID, "tweetId": l.TweetID})
	if err != nil {
		logger.Error(err)
	}

	return
}

func checkLiked(l *entity.Like) (liked bool) {
	follows, err := collection.Follows()
	if err != nil {
		logger.Error(err)
		return
	}
	defer follows.Close()

	n, err := follows.Find(bson.M{"userId": l.UserID, "tweetId": l.TweetID}).Count()
	if err != nil && err != mgo.ErrNotFound {
		logger.Error(err)
	}

	return n > 0
}

func checkValidLike(l *entity.Like) error {
	if !l.UserID.Valid() {
		return fmt.Errorf("invalid userId")
	}
	if !l.TweetID.Valid() {
		return fmt.Errorf("invalid tweetId")
	}
	return nil
}
