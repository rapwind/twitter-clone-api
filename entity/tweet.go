package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	// Tweet ... structure of a tweet
	Tweet struct {
		ID               bson.ObjectId `json:"id"                         bson:"_id"                        validate:"objectId"`
		Text             string        `json:"text"                       bson:"text"                       validate:"min=1"`
		UserID           bson.ObjectId `json:"-"                          bson:"userId"                     validate:"objectId"`
		InReplyToTweetID bson.ObjectId `json:"inReplyToTweetId,omitempty" bson:"inReplyToTweetId,omitempty" validate:"objectId"`
		CreatedAt        time.Time     `json:"createdAt"                  bson:"createdAt"`
		DeletedAt        *time.Time    `json:"-"                          bson:"deletedAt,omitempty"`
	}

	// TweetDetail ... structure of a tweet "more" information
	TweetDetail struct {
		*TweetDetailWithoutReply
		InReplyToTweet *TweetDetailWithoutReply `json:"inReplyToTweet,omitempty"`
	}

	// TweetDetailWithoutReply ... structure of a tweet "more" information
	TweetDetailWithoutReply struct {
		*Tweet
		User  *UserDetail `json:"user"`
		Liked *bool       `json:"liked,omitempty"`
	}
)
