package service

import (
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"gopkg.in/mgo.v2/bson"
)

// ReadUserByID gets "entity.User" data
func ReadUserByID(id bson.ObjectId) (u *entity.User, err error) {
	users, err := collection.Users()
	if err != nil {
		return
	}
	defer users.Close()

	u = new(entity.User)
	err = users.Find(bson.M{"_id": id}).One(u)
	return
}

// ReadFollowsCountsByID gets entity.UserDetail.FollowingCount and FollowerCount
func ReadFollowsCountsByID(id bson.ObjectId) (followingCount int, followerCount int, err error) {
	follows, err := collection.Follows()
	if err != nil {
		return
	}
	defer follows.Close()

	followingCount, err = follows.Find(bson.M{"userId": id}).Count()
	if err != nil {
		return
	}

	followerCount, err = follows.Find(bson.M{"targetId": id}).Count()
	return
}
