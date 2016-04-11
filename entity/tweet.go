package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Tweet ... structure of a tweet
type Tweet struct {
	ID               bson.ObjectId `json:"id"                         bson:"_id"                        validate:"objectId"`
	Text             string        `json:"text"                       bson:"text"                       validate:"min=1"`
	UserID           bson.ObjectId `json:"-"                          bson:"userId"                     validate:"objectId"`
	InReplyToUserID  bson.ObjectId `json:"-"                          bson:"inReplyToUserId,omitempty"  validate:"objectId"`
	InReplyToTweetID bson.ObjectId `json:"inReplyToTweetId,omitempty" bson:"inReplyToTweetId,omitempty" validate:"objectId"`
	CreatedAt        time.Time     `json:"createdAt"                  bson:"createdAt"`
	DeletedAt        *time.Time    `json:"-"                          bson:"deletedAt,omitempty"`
}
