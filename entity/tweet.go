package entity

import (
	"time"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Tweet ... structure of a tweet
	Tweet struct {
		ID                 bson.ObjectId `json:"id"                           bson:"_id"                          validate:"objectId"`
		Text               string        `json:"text"                         bson:"text"                         validate:"min=1"`
		ContentURL         string        `json:"contentUrl,omitempty"         bson:"contentUrl,omitempty"         validate:"min=1"`
		UserID             bson.ObjectId `json:"-"                            bson:"userId"                       validate:"objectId"`
		InReplyToTweetID   bson.ObjectId `json:"inReplyToTweetId,omitempty"   bson:"inReplyToTweetId,omitempty"   validate:"objectId"`
		InRetweetToTweetID bson.ObjectId `json:"inRetweetToTweetId,omitempty" bson:"inRetweetToTweetId,omitempty" validate:"objectId"`
		CreatedAt          time.Time     `json:"createdAt"                    bson:"createdAt"`
		DeletedAt          *time.Time    `json:"-"                            bson:"deletedAt,omitempty"`
	}

	// TweetDetail ... structure of a tweet "more" information
	TweetDetail struct {
		*TweetDetailWithoutReply
		InReplyToTweet   *TweetDetailWithoutReply `json:"inReplyToTweet,omitempty"`
		InRetweetToTweet *TweetDetailWithoutReply `json:"inRetweetToTweet,omitempty"`

		TargetFunc   func() int64  `json:"-"`
		PriorityFunc func() string `json:"-"`
	}

	// TweetDetailWithoutReply ... structure of a tweet "more" information
	TweetDetailWithoutReply struct {
		*Tweet
		User           *UserDetail `json:"user"`
		LikedCount     int         `json:"likedCount"`
		Liked          *bool       `json:"liked,omitempty"`
		RetweetedCount int         `json:"retweetedCount"`
		Retweeted      *bool       `json:"retweeted,omitempty"`
	}
)

// Target from Searcher interface
func (td *TweetDetail) Target() int64 {
	return td.TargetFunc()
}

// Priority from Searcher interface
func (td *TweetDetail) Priority() string {
	return td.PriorityFunc()
}

func initTweetsCollection() {

	// ensure index for tweets collection
	tweets, err := collection.Tweets()
	env.AssertErrForInit(err)

	defer tweets.Close()

	err = tweets.EnsureIndex(mgo.Index{
		Key:        []string{"-_id", "deletedAt"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)

	err = tweets.EnsureIndex(mgo.Index{
		Key:        []string{"userId", "deletedAt"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err

	err = tweets.EnsureIndex(mgo.Index{ // for retweetedCount, retweeted
		Key:        []string{"inRetweetToTweetId", "userId"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)
}
