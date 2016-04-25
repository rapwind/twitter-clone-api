package service

import (
	"fmt"

	"sync"

	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
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
func ReadUserDetailByID(id bson.ObjectId, l bson.ObjectId) (ud *entity.UserDetail, err error) {
	var wg sync.WaitGroup
	wg.Add(3)

	finChan := make(chan bool)
	errChan := make(chan error)
	userChan := make(chan *entity.User)
	followingChan := make(chan bool)
	countsChan := make(chan *entity.UserDetail)

	ud = new(entity.UserDetail)

	go func() {
		wg.Wait()
		finChan <- true
	}()

	go func(ud *entity.UserDetail) {
		defer wg.Done()

		u, err := ReadUserByID(id)
		if err != nil {
			errChan <- err
			return
		}
		userChan <- u
	}(ud)

	go func(id bson.ObjectId, l bson.ObjectId) {
		defer wg.Done()

		f := entity.Follow{
			UserID:   l,
			TargetID: id,
		}
		followingChan <- checkFollowing(&f)
	}(id, l)

	go func(id bson.ObjectId) {
		defer wg.Done()

		u := &entity.User{ID: id}
		d := &entity.UserDetail{User: u}
		fmt.Println(d)
		if err := AppendCounterToUserDetail(d); err != nil {
			errChan <- err
			return
		}
		countsChan <- d
	}(id)

LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case err = <-errChan:
			return
		case ud.User = <-userChan:
		case f := <-followingChan:
			ud.Following = &f
		case d := <-countsChan:
			ud.TweetsCount = d.TweetsCount
			ud.LikesCount = d.LikesCount
			ud.FollowerCount = d.FollowerCount
			ud.FollowingCount = d.FollowingCount
		}
	}
	return
}

// AppendCounterToUserDetail returns UserDetail append some count
func AppendCounterToUserDetail(ud *entity.UserDetail) (err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	finChan := make(chan bool)
	errChan := make(chan error)
	tweetsCountChan := make(chan int)
	likesCountChan := make(chan int)
	followingCountChan := make(chan int)
	followerCountChan := make(chan int)

	go func() {
		wg.Wait()
		finChan <- true
	}()

	go func(id bson.ObjectId) {
		defer wg.Done()
		tc, lc, err := readTweetsCountsByUserID(id)
		if err != nil {
			errChan <- err
			return
		}
		tweetsCountChan <- tc
		likesCountChan <- lc
	}(ud.ID)

	go func(id bson.ObjectId) {
		defer wg.Done()
		fgc, fwc, err := readFollowsCountsByID(id)
		if err != nil {
			errChan <- err
			return
		}
		followerCountChan <- fwc
		followingCountChan <- fgc
	}(ud.ID)

LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case err = <-errChan:
			return
		case ud.TweetsCount = <-tweetsCountChan:
		case ud.LikesCount = <-likesCountChan:
		case ud.FollowerCount = <-followerCountChan:
		case ud.FollowingCount = <-followingCountChan:
		}
	}
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

func readFollowsCountsByID(id bson.ObjectId) (followingCount int, followerCount int, err error) {
	var wg sync.WaitGroup
	wg.Add(2)

	finChan := make(chan bool)
	errChan := make(chan error)
	followingCountChan := make(chan int)
	followerCountChan := make(chan int)

	go func() {
		wg.Wait()
		finChan <- true
	}()

	go func(id bson.ObjectId) {
		defer wg.Done()

		follows, err := collection.Follows()
		if err != nil {
			return
		}
		defer follows.Close()

		c, err := follows.Find(bson.M{"userId": id}).Count()
		if err != nil {
			return
		}
		followingCountChan <- c
	}(id)

	go func(id bson.ObjectId) {
		defer wg.Done()

		follows, err := collection.Follows()
		if err != nil {
			errChan <- err
			return
		}
		defer follows.Close()

		c, err := follows.Find(bson.M{"targetId": id}).Count()
		if err != nil {
			errChan <- err
			return
		}
		followerCountChan <- c
	}(id)

LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case err = <-errChan:
			return
		case followingCount = <-followingCountChan:
		case followerCount = <-followerCountChan:
		}
	}

	return
}

// CreateFollow creates "entity.Follow" data
func CreateFollow(f *entity.Follow) (err error) {
	if err := checkRequiredForFollowing(f); err != nil {
		return errors.BadParams("follow", "invalid")
	}

	if checkFollowing(f) {
		err = errors.DataConflict()
		return
	}

	f.ID = bson.NewObjectId()
	f.CreatedAt = f.ID.Time()

	follows, err := collection.Follows()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer follows.Close()

	err = follows.Insert(f)
	if err != nil {
		logger.Error(err)
	}

	return
}

// RemoveFollow deletes a document on follow collection
func RemoveFollow(f *entity.Follow) (err error) {
	if err := checkRequiredForFollowing(f); err != nil {
		return errors.BadParams("follow", "invalid")
	}

	follows, err := collection.Follows()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer follows.Close()

	err = follows.Remove(bson.M{"userId": f.UserID, "targetId": f.TargetID})
	if err != nil && err != mgo.ErrNotFound {
		logger.Error(err)
	}

	return
}

// checkFollowing already checks a follow status
func checkFollowing(f *entity.Follow) (followed bool) {

	follows, err := collection.Follows()
	if err != nil {
		logger.Error(err)
		return
	}
	defer follows.Close()

	n, err := follows.Find(bson.M{"userId": f.UserID, "targetId": f.TargetID}).Count()
	if err != nil && err != mgo.ErrNotFound {
		logger.Error(err)
	}

	return n > 0
}

// checkRequiredForFollowing checks required fields of Follow
func checkRequiredForFollowing(f *entity.Follow) error {
	if !f.UserID.Valid() {
		return fmt.Errorf("invalid userId")
	}
	if !f.TargetID.Valid() {
		return fmt.Errorf("invalid targetId")
	}
	if f.UserID == f.TargetID {
		return fmt.Errorf("invalid targetId")
	}
	return nil
}
