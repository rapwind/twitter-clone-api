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

// CreateInstallation inserts a installation
func CreateInstallation(i *entity.Installation) (err error) {

	if i.ID != "" {
		return fmt.Errorf("already objectId, oid = %s", i.ID)
	}

	i.ID = bson.NewObjectId()
	i.UUID = utils.GetNewUUIDv4()
	i.CreatedAt = i.ID.Time()
	i.UpdatedAt = i.ID.Time()

	installations, err := collection.Installations()
	if err != nil {
		return
	}
	defer installations.Close()

	err = installations.Insert(i)
	if err != nil && !mgo.IsDup(err) {
		logger.Error(err)
	}

	return
}

// UpdateInstallation inserts a installation
func UpdateInstallation(i *entity.Installation) (err error) {

	i.UpdatedAt = utils.GetNowTruncateSecond()

	installations, err := collection.Installations()
	if err != nil {
		return
	}
	defer installations.Close()

	if err = installations.UpdateId(i.ID, i); err != nil {
		logger.Error(err)
		return
	}

	return
}

// ReadInstallationByUUID gets "entity.Installation" data
func ReadInstallationByUUID(uuid string) (i *entity.Installation, err error) {
	installations, err := collection.Installations()
	if err != nil {
		return
	}
	defer installations.Close()

	i = new(entity.Installation)
	err = installations.Find(bson.M{"uuid": uuid}).One(i)
	return
}
