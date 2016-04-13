package service

import (
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ReadSessionByID gets "entity.Session" data
func ReadSessionByID(id bson.ObjectId) (s *entity.Session, err error) {
	sessions, err := collection.Sessions()
	if err != nil {
		return
	}
	defer sessions.Close()

	s = new(entity.Session)
	err = sessions.Find(bson.M{"_id": id}).One(s)
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
