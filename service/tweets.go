package service

import (
	"fmt"

	"sync"

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

// ReadLatestTweet returns the newest tweet of a user.
func ReadLatestTweet(userID bson.ObjectId) (t *entity.Tweet, err error) {
	tweets, err := collection.Tweets()
	if err != nil {
		logger.Error(err)
		return
	}
	defer tweets.Close()

	t = &entity.Tweet{}
	err = tweets.Find(bson.M{"userId": userID, "deletedAt": bson.M{"$exists": false}}).Sort("-_id").One(t)
	return
}

// ReadUserLikedTweetDetails returns TweetDetails by user ID
func ReadUserLikedTweetDetails(userID bson.ObjectId, loginUserID bson.ObjectId, limit int, maxID bson.ObjectId) (tds []*entity.TweetDetail, err error) {
	if limit == 0 {
		return
	}

	m := bson.M{"userId": userID}

	// if maxID is set:
	if maxID.Valid() {
		m["tweetId"] = bson.M{"$lte": maxID}
	}

	likes, err := collection.Likes()
	if err != nil {
		return
	}
	defer likes.Close()

	iter := likes.Find(m).Sort("-createdAt").Iter()

	ts := []entity.Tweet{}

	var l entity.Like
	var t *entity.Tweet
	for iter.Next(&l) {
		t, err = ReadTweetByID(l.TweetID)
		if err != nil {
			return
		}

		if t.DeletedAt == nil { // removed tweets are not appended
			ts = append(ts, *t)
			if limit > 0 && len(ts) == limit {
				break
			}
		}
	}

	err = iter.Close()
	if err != nil {
		return
	}

	tds, err = readTweetDetailsByTweets(ts, loginUserID)

	return
}

// ReadUserTweetDetails returns TweetDetails by user ID
func ReadUserTweetDetails(userID bson.ObjectId, loginUserID bson.ObjectId, limit int, maxID bson.ObjectId) (tds []*entity.TweetDetail, err error) {
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

// ReadTweetDetails returns an array of TweetDetail(s)
func ReadTweetDetails(limit int, maxID bson.ObjectId, sinceID bson.ObjectId, userID bson.ObjectId, following bool, q string) (tds []*entity.TweetDetail, err error) {
	m := []bson.M{
		bson.M{"deletedAt": bson.M{"$exists": false}},
	}

	// if maxId is set:
	if maxID.Valid() {
		m = append(m, bson.M{"_id": bson.M{"$lte": maxID}})
	}

	// if sinceId is set:
	if sinceID.Valid() {
		m = append(m, bson.M{"_id": bson.M{"$gt": sinceID}})
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

		if len(flws) == 0 {
			// The login user follows no users.
			m = append(m, bson.M{"userId": userID})
		} else {
			// Convert an array of Follow(s) into an array of user IDs
			n := len(flws)
			ids := make([]bson.ObjectId, n+1)
			for i, v := range flws {
				ids[i] = v.TargetID
			}
			ids[n] = userID
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

func readSortedTweetDetails(m bson.M, limit int, loginUserID bson.ObjectId) (tds []*entity.TweetDetail, err error) {
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

	tds, err = readTweetDetailsByTweets(ts, loginUserID)
	return
}

func readTweetDetailsByTweets(ts []entity.Tweet, loginUserID bson.ObjectId) (tds []*entity.TweetDetail, err error) {
	var wg sync.WaitGroup

	finChan := make(chan bool)
	tweetsChan := make(chan *entity.TweetDetail, len(ts))

	wg.Add(len(ts))
	go func() {
		wg.Wait()
		finChan <- true
	}()

	for _, t := range ts {
		go func(t entity.Tweet) {
			defer wg.Done()

			td := new(entity.TweetDetail)
			td, err = readTweetDetailByTweet(t, loginUserID)
			if err != nil {
				if err != mgo.ErrNotFound {
					logger.Error(err)
				}
				return
			}
			td.TargetFunc = func() int64 { return td.CreatedAt.Unix() }
			td.PriorityFunc = func() string { return td.ID.Hex() }

			tweetsChan <- td
		}(t)
	}
	s := &entity.SortedSlice{DESC: true}
LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case pd := <-tweetsChan:
			s.Add(pd)
		}
	}

	for _, i := range s.S {
		tds = append(tds, i.(*entity.TweetDetail))
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

	td = &entity.TweetDetail{TweetDetailWithoutReply: tdwr, InReplyToTweet: inReplyToTweet}
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
	var wg sync.WaitGroup
	wg.Add(3)

	finChan := make(chan bool)
	errChan := make(chan error)
	userDetailChan := make(chan *entity.UserDetail)
	likedChan := make(chan *bool)
	likedCountChan := make(chan int)

	tdwr = &entity.TweetDetailWithoutReply{Tweet: &t}

	go func() {
		wg.Wait()
		finChan <- true
	}()

	go func(id bson.ObjectId) {
		defer wg.Done()

		u, err := ReadUserDetailByID(id, loginUserID)
		if err != nil {
			errChan <- err
			return
		}
		userDetailChan <- u
	}(t.UserID)

	go func(loginUserID bson.ObjectId, tweetID bson.ObjectId) {
		defer wg.Done()

		likes, err := collection.Likes()
		if err != nil {
			errChan <- err
			return
		}
		defer likes.Close()
		l := entity.Like{}
		var liked bool
		if loginUserID.Valid() && likes.Find(bson.M{"userId": loginUserID, "tweetId": tweetID}).One(&l) == nil {
			liked = true
		} else {
			liked = false
		}
		likedChan <- &liked
	}(loginUserID, t.ID)

	go func(id bson.ObjectId) {
		defer wg.Done()
		c, err := ReadLikedCountByTweetID(id)
		if err != nil {
			errChan <- err
			return
		}
		likedCountChan <- c
	}(t.ID)

LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case err = <-errChan:
			return
		case tdwr.User = <-userDetailChan:
		case tdwr.Liked = <-likedChan:
		case tdwr.LikedCount = <-likedCountChan:
		}
	}
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

// ReadLikedCountByTweetID gets entity.TweetDetail.LikedCount by TweetID
func ReadLikedCountByTweetID(id bson.ObjectId) (likedCount int, err error) {
	likes, err := collection.Likes()
	if err != nil {
		return
	}
	defer likes.Close()

	likedCount, err = likes.Find(bson.M{"tweetId": id}).Count()

	return
}

func readTweetsCountsByUserID(uid bson.ObjectId) (tweetsCount int, likesCount int, err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	finChan := make(chan bool)
	errChan := make(chan error)
	tweetsCountChan := make(chan int)
	likesCountChan := make(chan int)

	go func() {
		wg.Wait()
		finChan <- true
	}()

	go func(id bson.ObjectId) {
		defer wg.Done()

		tweets, err := collection.Tweets()
		if err != nil {
			errChan <- err
			return
		}
		defer tweets.Close()

		c, err := tweets.Find(bson.M{"userId": id}).Count()
		if err != nil {
			errChan <- err
			return
		}
		tweetsCountChan <- c
	}(uid)

	go func(id bson.ObjectId) {
		defer wg.Done()

		likes, err := collection.Likes()
		if err != nil {
			errChan <- err
			return
		}
		defer likes.Close()

		c, err := likes.Find(bson.M{"userId": id}).Count()
		if err != nil {
			errChan <- err
			return
		}
		likesCountChan <- c
	}(uid)

LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case err = <-errChan:
			return
		case tweetsCount = <-tweetsCountChan:
		case likesCount = <-likesCountChan:
		}
	}

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
