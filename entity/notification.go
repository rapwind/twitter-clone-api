package entity

import (
	"time"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// PushMessage ... structure of a push message
	PushMessage struct {
		ID    bson.ObjectId
		Text  string
		Count int
		Type  string
	}

	// CommonNotification ... fields commonly used in all types of notifications
	CommonNotification struct {
		ID        bson.ObjectId `json:"id"        bson:"_id"    validate:"objectId"`
		UserID    bson.ObjectId `json:"-"         bson:"userId" validate:"objectId"`
		Type      string        `json:"type"      bson:"-"`
		Unread    bool          `json:"unread"    bson:"unread"`
		CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	}

	// FollowNotification ... fields only used in follow notifications (MongoDB)
	FollowNotification struct {
		UserID bson.ObjectId `json:"-" bson:"userId" validate:"objectId"`
	}

	// ReplyNotification ... fields only used in reply notifications (MongoDB)
	ReplyNotification struct {
		TweetID bson.ObjectId `json:"-" bson:"tweetId" validate:"objectId"`
	}

	// LikeNotification ... fields only used in like notifications (MongoDB)
	LikeNotification struct {
		UserID  bson.ObjectId `json:"-" bson:"userId"  validate:"objectId"`
		TweetID bson.ObjectId `json:"-" bson:"tweetId" validate:"objectId"`
	}

	// Notification ... structure of a notification (MongoDB)
	Notification struct {
		CommonNotification `bson:",inline"`
		Follow             *FollowNotification `json:"-" bson:"follow,omitempty"`
		Reply              *ReplyNotification  `json:"-" bson:"reply,omitempty"`
		Like               *LikeNotification   `json:"-" bson:"like,omitempty"`
	}

	// FollowNotificationDetail ... fields only used in follow notifications (API response)
	FollowNotificationDetail struct {
		User UserDetail `json:"user"`
	}

	// ReplyNotificationDetail ... fields only used in reply notifications (API Response)
	ReplyNotificationDetail struct {
		Tweet TweetDetail `json:"tweet"`
	}

	// LikeNotificationDetail ... fields only used in like notifications (API response)
	LikeNotificationDetail struct {
		User  UserDetail  `json:"user"`
		Tweet TweetDetail `json:"tweet"`
	}

	// NotificationDetail ... structure of a notification (API response)
	NotificationDetail struct {
		CommonNotification
		Follow *FollowNotificationDetail `json:"follow,omitempty"`
		Reply  *ReplyNotificationDetail  `json:"reply,omitempty"`
		Like   *LikeNotificationDetail   `json:"like,omitempty"`
	}
)

func initNotificationsCollection() {
	likes, err := collection.Notifications()
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
}
