package service

import (
	"fmt"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/logger"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CreateSession creates "entity.Session" data
func CreateSession(s *entity.Session) (err error) {
	if s.ID != "" {
		return fmt.Errorf("already objectId, oid = %s", s.ID)
	}

	sessions, err := collection.Sessions()
	if err != nil {
		return
	}
	defer sessions.Close()

	err = sessions.Find(bson.M{"userId": s.UserID, "installationId": s.InstallationID}).One(s)
	if s.ID != "" {
		RemoveSessionByID(s.ID)
	}

	s.ID = bson.NewObjectId()
	s.UUID = utils.GetNewUUIDv4()
	s.CreatedAt = s.ID.Time()

	err = sessions.Insert(s)
	if err != nil && !mgo.IsDup(err) {
		logger.Error(err)
	}

	return
}

// ReadSessionByUUID gets "entity.Session" data
func ReadSessionByUUID(uuid string) (s *entity.Session, err error) {
	sessions, err := collection.Sessions()
	if err != nil {
		return
	}
	defer sessions.Close()

	s = new(entity.Session)
	err = sessions.Find(bson.M{"uuid": uuid}).One(s)
	return
}

// ReadSessionsByUser gets "entity.Session" list
func ReadSessionsByUser(user *entity.User) (ss []entity.Session, err error) {
	sessions, err := collection.Sessions()
	if err != nil {
		return
	}
	defer sessions.Close()

	err = sessions.Find(bson.M{"userId": user.ID}).All(ss)
	return
}

// RemoveSessionByID deletes a document on follow collection
func RemoveSessionByID(id bson.ObjectId) (err error) {
	sessions, err := collection.Sessions()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer sessions.Close()

	err = sessions.Remove(bson.M{"_id": id})
	if err != nil && err != mgo.ErrNotFound {
		logger.Error(err)
	}

	return
}

// RemoveSessionByUUID deletes a document on follow collection
func RemoveSessionByUUID(uuid string) (err error) {
	sessions, err := collection.Sessions()
	if err != nil {
		logger.Error(err)
		return err
	}
	defer sessions.Close()

	err = sessions.Remove(bson.M{"uuid": uuid})
	if err != nil && err != mgo.ErrNotFound {
		logger.Error(err)
	}

	return
}
