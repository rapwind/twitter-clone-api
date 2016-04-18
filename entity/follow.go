package entity

import (
	"time"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Follow ... structure of a follow
	Follow struct {
		ID        bson.ObjectId `json:"-"        bson:"_id"      validate:"objectId"`
		UserID    bson.ObjectId `json:"userId"   bson:"userId"   validate:"objectId"`
		TargetID  bson.ObjectId `json:"targetId" bson:"targetId" validate:"objectId"`
		CreatedAt time.Time     `json:"-"        bson:"createdAt"`
	}
)

func initFollowsCollection() {
	follows, err := collection.Follows()
	env.AssertErrForInit(err)

	defer follows.Close()

	err = follows.EnsureIndex(mgo.Index{
		Key:        []string{"userId", "-createdAt"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)

	err = follows.EnsureIndex(mgo.Index{
		Key:        []string{"targetId", "-createdAt"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)

	err = follows.EnsureIndex(mgo.Index{
		Key:        []string{"userId", "targetId"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)
}
