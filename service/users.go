package service

import (
	"fmt"

	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/logger"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CreateUser creates "entity.User" data
func CreateUser(u *entity.User) (err error) {
	if u.ID != "" {
		return fmt.Errorf("already objectId, oid = %s", u.ID)
	}

	u.ID = bson.NewObjectId()
	u.CreatedAt = u.ID.Time()
	u.UpdatedAt = u.ID.Time()

	if u.ScreenName == "" {
		u.ScreenName = utils.RandString(constant.DefaultScreenNameSize)
	}

	users, err := collection.Users()
	if err != nil {
		return
	}
	defer users.Close()

	err = users.Insert(u)
	if err != nil && !mgo.IsDup(err) {
		logger.Error(err)
	}

	return
}

// ReadUserDetailByID returns UserDetail by user ID
func ReadUserDetailByID(id bson.ObjectId) (ud *entity.UserDetail, err error) {
	u, err := ReadUserByID(id)
	if err != nil {
		return
	}

	tweetsCount, likesCount, err := ReadTweetsCountsByUser(*u)
	if err != nil {
		return
	}

	followingCount, followerCount, err := ReadFollowsCountsByID(id)
	if err != nil {
		return
	}

	ud = &entity.UserDetail{u, tweetsCount, likesCount, followerCount, followingCount, nil}
	return
}

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

// ReadFollowingByID returns an array of entity.Follow
func ReadFollowingByID(id bson.ObjectId, offset int, limit int) (flws []entity.Follow, err error) {
	follows, err := collection.Follows()
	if err != nil {
		return
	}
	defer follows.Close()

	q := follows.Find(bson.M{"userId": id}).Sort("-_id")

	if offset > 0 {
		q = q.Skip(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	flws = []entity.Follow{}
	err = q.All(&flws)
	return
}

// ReadFollowerByID returns an array of entity.Follow
func ReadFollowerByID(id bson.ObjectId, offset int, limit int) (flws []entity.Follow, err error) {
	follows, err := collection.Follows()
	if err != nil {
		return
	}
	defer follows.Close()

	q := follows.Find(bson.M{"targetId": id}).Sort("-_id")

	if offset > 0 {
		q = q.Skip(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	flws = []entity.Follow{}
	err = q.All(&flws)
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
