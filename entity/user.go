package entity

import (
	"time"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// User ... structure of a user
	User struct {
		ID                        bson.ObjectId `json:"id"                                  bson:"_id"                                 validate:"objectId"`
		Name                      string        `json:"name"                                bson:"name"                                validate:"min=1"`
		ScreenName                string        `json:"screenName"                          bson:"screenName"                          validate:"min=1"`
		Email                     string        `json:"-"                                   bson:"email"                               validate:"min=1"`
		PhoneNumber               string        `json:"-"                                   bson:"phoneNumber"                         validate:"min=1"`
		PasswordHash              string        `json:"-"                                   bson:"passwordHash"                        validate:"min=1"`
		ProfileImageURL           string        `json:"profileImageUrl,omitempty"           bson:"profileImageUrl,omitempty"`
		ProfileBackgroundImageURL string        `json:"profileBackgroundImageUrl,omitempty" bson:"profileBackgroundImageUrl,omitempty"`
		Biography                 string        `json:"biography,omitempty"                 bson:"biography,omitempty"`
		LocationText              string        `json:"locationText,omitempty"              bson:"locationText,omitempty"`
		URL                       string        `json:"url,omitempty"                       bson:"url,omitempty"`
		Birthday                  *time.Time    `json:"birthday,omitempty"                  bson:"birthday,omitempty"`
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

		TargetFunc   func() int64  `json:"-"`
		PriorityFunc func() string `json:"-"`
	}

	// UserRegisterRequest ... structure of a user register request
	UserRegisterRequest struct {
		Name         string `json:"name"         bson:"name"         validate:"min=1"`
		ScreenName   string `json:"screenName"   bson:"screenName"   validate:"min=1"`
		PhoneNumber  string `json:"phoneNumber"  bson:"phoneNumber"  validate:"min=1"`
		Email        string `json:"email"        bson:"email"        validate:"min=1"`
		PasswordHash string `json:"passwordHash" bson:"passwordHash" validate:"min=1"`
	}
)

// Target from Searcher interface
func (ud *UserDetail) Target() int64 {
	return ud.TargetFunc()
}

// Priority from Searcher interface
func (ud *UserDetail) Priority() string {
	return ud.PriorityFunc()
}

func initUsersCollection() {

	// ensure index for users collection
	users, err := collection.Users()
	env.AssertErrForInit(err)

	defer users.Close()

	err = users.EnsureIndex(mgo.Index{
		Key:        []string{"screenName"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)

	err = users.EnsureIndex(mgo.Index{
		Key:        []string{"phoneNumber"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)
}
