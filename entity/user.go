package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	// User ... structure of a user
	User struct {
		UserID                    bson.ObjectId `json:"userId"                              bson:"userId"                              validate:"objectId"`
		Name                      string        `json:"name"                                bson:"name"                                validate:"min=1"`
		ScreenName                string        `json:"screenName"                          bson:"screenName"                          validate:"min=1"`
		ProfileImageURL           string        `json:"profileImageUrl,omitempty"           bson:"profileImageUrl,omitempty"`
		ProfileBackgroundImageURL string        `json:"profileBackgroundImageUrl,omitempty" bson:"profileBackgroundImageUrl,omitempty"`
		Biography                 string        `json:"biography,omitempty"                 bson:"biography,omitempty"`
		LocationText              string        `json:"locationText,omitempty"              bson:"localtionText,omitempty"`
		URL                       string        `json:"url,omitempty"                       bson:"url,omitempty"`
		Birthday                  *time.Time    `json:"time,omitempty"                      bson:"time,omitempty"`
		CreatedAt                 time.Time     `json:"createdAt"                           bson:"createdAt"`
		UpdatedAt                 time.Time     `json:"updatedAt"                           bson:"updatedAt"`
	}

	// UserDetail ... structure of a user "more" information
	UserDetail struct {
		*User
		TweetsCount    int   `json:"tweetsCount"`
		LikesCount     int   `json:"likesCount"`
		FollowerCount  int   `json:"followerCount"`
		FollowingCount int   `json:"followingCount"`
		Following      *bool `json:"following,omitempty"`
	}
)
