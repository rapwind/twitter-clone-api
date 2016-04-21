package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	// Notification ... structure of a notification
	Notification struct {
		ID        bson.ObjectId `json:"id"        bson:"_id"  validate:"objectId"`
		Text      string        `json:"text"      bson:"text" validate:"min=1"`
		Read      bool          `json:"read"      bson:"read"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
		DeletedAt *time.Time    `json:"-"         bson:"deletedAt,omitempty"`
	}
)
