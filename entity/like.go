package entity

import (
	"time"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (

	// Like ... structure of a like
	Like struct {
		ID        bson.ObjectId `json:"-"       bson:"_id"     validate:"objectId"`
		UserID    bson.ObjectId `json:"userId"  bson:"userId"  validate:"objectId"`
		TweetID   bson.ObjectId `json:"tweetId" bson:"tweetId" validate:"objectId"`
		CreatedAt time.Time     `json:"-"       bson:"createdAt"`
	}
)

func initLikesCollection() {
	likes, err := collection.Likes()
	env.AssertErrForInit(err)

	defer likes.Close()

	err = likes.EnsureIndex(mgo.Index{
		Key:        []string{"userId"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)

	err = likes.EnsureIndex(mgo.Index{
		Key:        []string{"tweetId"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)

	err = likes.EnsureIndex(mgo.Index{
		Key:        []string{"userId", "tweetId"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)
}
