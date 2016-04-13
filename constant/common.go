package constant

import (
	"github.com/techcampman/twitter-d-server/entity"
	"gopkg.in/mgo.v2/bson"
)

type (
	// FollowsKey ... constant.FollowingIDKey or constant.FollowerIDKey
	FollowsKey string
)

const (
	// FollowingIDKey ... the key for a following user ID (followee)
	FollowingIDKey = "userId"

	// FollowerIDKey ... the key for a followed user ID
	FollowerIDKey = "targetId"
)

// GetFollowingIDKey returns a following user ID from entity.Follow.
func GetFollowingIDKey (v entity.Follow) bson.ObjectId {
	return v.UserID
}

// GetFollowerIDKey returns a followed user ID from entity.Follow.
func GetFollowerIDKey (v entity.Follow) bson.ObjectId {
	return v.TargetID
}