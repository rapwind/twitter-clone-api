package entity

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	// Notification ... structure of a notification
	Notification struct {
		ID               bson.ObjectId `json:"id"        bson:"_id"  validate:"objectId"`
		Text             string        `json:"text"      bson:"text" validate:"min=1"`
		Read             bool          `json:"read"      bson:"read"`
		CreatedAt        time.Time     `json:"createdAt" bson:"createdAt"`
		DeletedAt        *time.Time    `json:"-"         bson:"deletedAt,omitempty"`
	}
)
