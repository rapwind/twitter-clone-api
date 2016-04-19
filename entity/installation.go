package entity

import (
	"time"

	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/env"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	// Installation ... structure of a installation
	Installation struct {
		ID          bson.ObjectId `json:"-"                     bson:"_id" validate:"objectId"`
		UUID        string        `json:"id"                    bson:"uuid"`
		ClientType  string        `json:"clientType"            bson:"clientType"`
		DeviceToken string        `json:"deviceToken,omitempty" bson:"deviceToken,omitempty"`
		CreatedAt   time.Time     `json:"createdAt"             bson:"createdAt"`
		UpdatedAt   time.Time     `json:"updatedAt"             bson:"updatedAt"`
	}
)

func initInstallationsCollection() {

	// ensure index for users collection
	installations, err := collection.Installations()
	env.AssertErrForInit(err)

	defer installations.Close()

	err = installations.EnsureIndex(mgo.Index{
		Key:        []string{"uuid"},
		Unique:     true,
		DropDups:   false,
		Background: true,
		Sparse:     true,
	})
	env.AssertErrForInit(err)
}
