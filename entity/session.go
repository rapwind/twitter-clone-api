package entity

import (
	"time"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Session ... structure of a session
	Session struct {
		ID             bson.ObjectId `bson:"_id"            validate:"objectId"`
		UUID           string        `bson:"uuid"`
		UserID         bson.ObjectId `bson:"userId"         validate:"objectId"`
		InstallationID bson.ObjectId `bson:"installationId" validate:"objectId"`
		CreatedAt      time.Time     `bson:"createdAt"`
	}

	// SessionRequest ... structure of a session request
	SessionRequest struct {
		AccountName  string `json:"accountName"`
		PasswordHash string `json:"passwordHash"`
	}
)

func initSessionsCollection() {

	// ensure index for users collection
	sessions, err := collection.Sessions()
	env.AssertErrForInit(err)

	defer sessions.Close()

	err = sessions.EnsureIndex(mgo.Index{
		Key:        []string{"uuid"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)
}
