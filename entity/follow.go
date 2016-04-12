package entity

import "gopkg.in/mgo.v2/bson"

type (
	// Follow ... structure of a follow
	Follow struct {
		ID       bson.ObjectId `json:"-"        bson:"_id"      validate:"objectId"`
		UserID   bson.ObjectId `json:"userId"   bson:"userId"   validate:"objectId"`
		TargetID bson.ObjectId `json:"targetId" bson:"targetId" validate:"objectId"`
	}
)
