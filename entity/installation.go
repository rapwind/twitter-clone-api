package entity

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type (
	// Installation ... structure of a installation
	Installation struct {
		ID          bson.ObjectId `json:"-"                     bson:"_id" validate:"objectId"`
		ClientType  string        `json:"clientType"            bson:"clientType"`
		DeviceToken string        `json:"deviceToken,omitempty" bson:"deviceToken,omitempty"`
		CreatedAt   time.Time     `json:"createdAt"             bson:"createdAt"`
	}

	// InstallationHeader ... structure of a installation header
	InstallationHeader struct {
		ID bson.ObjectId `json:"-"`
	}
)
