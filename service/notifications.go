package service

import (
	"fmt"
	"sync"

	"errors"

	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/env"
	"github.com/techcampman/twitter-d-server/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CreateFollowNotification creates a follow notification.
func CreateFollowNotification(f *entity.Follow) (err error) {
	// The followee (f.TargetID) receives a follow notification from the follower (f.UserID).
	n := &entity.Notification{
		Follow: &entity.FollowNotification{UserID: f.UserID},
	}
	recvUserID := f.TargetID
	createNotification(recvUserID, n)

	// Send a follow notification.
	u, err := ReadUserByID(recvUserID)
	if err != nil {
		return
	}
	pm := &entity.PushMessage{
		ID:   n.ID,
		Type: constant.NotificationTypeFollow,
		Text: fmt.Sprintf("@%sさんがフォローしました", u.ScreenName),
	}
	sendNotificationForUser(u, pm)

	return
}

// CreateReplyNotification creates a reply notification.
func CreateReplyNotification(t *entity.Tweet) (err error) {
	t0, err := ReadTweetByID(t.InReplyToTweetID)
	if err != nil {
		return
	}

	// The writer of the original tweet (t0.UserID) receives a reply notification from the user that replies (t.UserID).
	n := &entity.Notification{
		Reply: &entity.ReplyNotification{TweetID: t.ID},
	}
	recvUserID := t0.UserID
	createNotification(recvUserID, n)

	// Send a reply notification.
	u, err := ReadUserByID(recvUserID)
	if err != nil {
		return
	}
	pm := &entity.PushMessage{
		ID:   n.ID,
		Type: constant.NotificationTypeReply,
		Text: fmt.Sprintf("@%sさんが@ツイートしました\n%s", u.ScreenName, t.Text),
	}
	sendNotificationForUser(u, pm)

	return
}

// CreateLikeNotification creates a reply notification.
func CreateLikeNotification(l *entity.Like) (err error) {
	t0, err := ReadTweetByID(l.TweetID)
	if err != nil {
		return
	}

	// The writer of the tweet (t0.TweetID) receives a like notification from the user that likes the tweet (l.UserID).
	n := &entity.Notification{
		Like: &entity.LikeNotification{UserID: l.UserID, TweetID: t0.ID},
	}
	recvUserID := t0.UserID
	createNotification(recvUserID, n)

	// Send a like notification.
	u, err := ReadUserByID(recvUserID)
	if err != nil {
		return
	}
	pm := &entity.PushMessage{
		ID:   n.ID,
		Type: constant.NotificationTypeLike,
		Text: fmt.Sprintf("@%sさんがいいねしました", u.ScreenName),
	}
	sendNotificationForUser(u, pm)

	return
}

func createNotification(userID bson.ObjectId, n *entity.Notification) (err error) {
	n.ID = bson.NewObjectId()
	n.CreatedAt = n.ID.Time()
	n.Unread = true
	n.UserID = userID

	nots, err := collection.Notifications()
	if err != nil {
		return
	}
	defer nots.Close()

	err = nots.Insert(n)
	if err != nil && !mgo.IsDup(err) {
		logger.Error(err)
	}
	return
}

// ReadNotificationDetails returns NotificationDetail(s)
func ReadNotificationDetails(userID bson.ObjectId, limit int, maxID bson.ObjectId, sinceID bson.ObjectId) (nds []entity.NotificationDetail, err error) {
	ns, err := readNotifications(userID, limit, maxID, sinceID)
	if err != nil {
		return
	}

	// Convert Notification(s) into NotificationDetail(s)
	nds = make([]entity.NotificationDetail, len(ns))
	var nd *entity.NotificationDetail
	for i, n := range ns {
		nd = &entity.NotificationDetail{n.CommonNotification, nil, nil, nil}

		if n.Follow != nil { // This is a follow notification
			nd.Type = constant.NotificationTypeFollow
			nd.Follow, err = readFollowNotificationDetail(*n.Follow, userID)
		} else if n.Reply != nil { // This is a reply notification.
			nd.Type = constant.NotificationTypeReply
			nd.Reply, err = readReplyNotificationDetail(*n.Reply, userID)
		} else if n.Like != nil { // This is a like notification.
			nd.Type = constant.NotificationTypeLike
			nd.Like, err = readLikeNotificationDetail(*n.Like, userID)
		} else {
			err = errors.New("Illegal notification type")
		}

		if err != nil {
			return
		}

		nds[i] = *nd
	}

	return
}

func readLikeNotificationDetail(n entity.LikeNotification, loginUserID bson.ObjectId) (nd *entity.LikeNotificationDetail, err error) {
	u, err := ReadUserDetailByID(n.UserID, loginUserID)
	if err != nil {
		return
	}
	t, err := ReadTweetDetailByID(n.TweetID, loginUserID)
	if err != nil {
		return
	}
	nd = &entity.LikeNotificationDetail{User: *u, Tweet: *t}
	return
}

func readReplyNotificationDetail(n entity.ReplyNotification, loginUserID bson.ObjectId) (nd *entity.ReplyNotificationDetail, err error) {
	t, err := ReadTweetDetailByID(n.TweetID, loginUserID)
	if err != nil {
		return
	}
	nd = &entity.ReplyNotificationDetail{Tweet: *t}
	return
}

func readFollowNotificationDetail(n entity.FollowNotification, loginUserID bson.ObjectId) (nd *entity.FollowNotificationDetail, err error) {
	u, err := ReadUserDetailByID(n.UserID, loginUserID)
	if err != nil {
		return
	}
	nd = &entity.FollowNotificationDetail{User: *u}
	return
}

func readNotifications(userID bson.ObjectId, limit int, maxID bson.ObjectId, sinceID bson.ObjectId) (ns []entity.Notification, err error) {
	m := []bson.M{
		bson.M{"userId": userID},
	}

	// if maxId is set:
	if maxID.Valid() {
		m = append(m, bson.M{"_id": bson.M{"$lte": maxID}})
	}

	// if sinceId is set:
	if sinceID.Valid() {
		m = append(m, bson.M{"_id": bson.M{"$gt": sinceID}})
	}

	nots, err := collection.Notifications()
	if err != nil {
		return
	}
	defer nots.Close()

	ns = []entity.Notification{}
	err = nots.Find(bson.M{"$and": m}).Sort("-createdAt").Limit(limit).All(&ns)

	return
}

// UpdateNotificationsUnread makes "unread" flag "false" in notifications from "minID" to "maxID"
func UpdateNotificationsUnread(userID bson.ObjectId, maxID bson.ObjectId, minID bson.ObjectId) (err error) {
	m := bson.M{"$and": []bson.M{
		bson.M{"userId": userID},
		bson.M{"_id": bson.M{"$lte": maxID}},
		bson.M{"_id": bson.M{"$gte": minID}},
	},
	}

	nots, err := collection.Notifications()
	if err != nil {
		return
	}

	_, err = nots.UpdateAll(m, bson.M{"$set": bson.M{"unread": false}})
	return
}

// ReadNotificationsCount counts the number of unread notifications.
func ReadNotificationsCount(userID bson.ObjectId) (n int, err error) {
	nots, err := collection.Notifications()
	if err != nil {
		return
	}

	n, err = nots.Find(bson.M{"userId": userID, "unread": true}).Count()
	return
}

func sendNotificationForUser(u *entity.User, pm *entity.PushMessage) (err error) {

	pm.Count, err = ReadNotificationsCount(u.ID)
	if err != nil {
		return
	}

	ss, err := ReadSessionsByUser(u)

	var wg sync.WaitGroup
	finChan := make(chan bool)
	installationsChan := make(chan *entity.Installation, len(ss))

	wg.Add(len(ss))
	go func() {
		wg.Wait()
		finChan <- true
	}()

	for _, s := range ss {
		go func(s entity.Session) {
			defer wg.Done()

			i, err := ReadInstallationByID(s.InstallationID)
			if err != nil {
				if err != mgo.ErrNotFound {
					logger.Error(err)
				}
				return
			}
			installationsChan <- i
		}(s)
	}
LOOP:
	for {
		select {
		case <-finChan:
			break LOOP
		case i := <-installationsChan:
			env.GetPushMessage().Send(pm.Text, pm.Count, pm.Type, pm.ID, i.ClientType, i.ArnEndpoint)
		}
	}

	fmt.Println(ss)

	return
}
