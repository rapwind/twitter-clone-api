package service

import (
	"fmt"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/logger"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CreateInstallation inserts a installation
func CreateInstallation(i *entity.Installation) (err error) {

	if i.ID != "" {
		return fmt.Errorf("already objectId, oid = %s", i.ID)
	}

	i.ID = bson.NewObjectId()
	i.CreatedAt = i.ID.Time()

	c, err := collection.Installations()
	if err != nil {
		return
	}
	defer c.Close()

	err = c.Insert(i)
	if err != nil && !mgo.IsDup(err) {
		logger.Error(err)
	}

	return
}
