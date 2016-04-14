package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	// Session ... structure of a session
	Session struct {
		ID             bson.ObjectId `bson:"_id"            validate:"objectId"`
		UserID         bson.ObjectId `bson:"userId"         validate:"objectId"`
		InstallationID bson.ObjectId `bson:"installationId" validate:"objectId"`
		CreatedAt      time.Time     `bson:"createdAt"`
	}

	// SessionRequest ... structure of a session request
	SessionRequest struct {
		ScreenName   string `json:"screenName"`
		PasswordHash string `json:"passwordHash"`
	}
)
