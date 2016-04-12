package entity

import "gopkg.in/mgo.v2/bson"

type (
	// Follow ... structure of a follow
	Follow struct {
		ID       bson.ObjectId `json:"-"        bson:"_id"`
		UserID   bson.ObjectId `json:"userId"   bson:"userId"`
		TargetID bson.ObjectId `json:"targetId" bson:"targetId"`
	}
)