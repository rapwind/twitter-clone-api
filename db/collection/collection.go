package collection

import (
	"github.com/techcampman/twitter-d-server/db/mongo"
	"github.com/techcampman/twitter-d-server/env"
)

const (
	installations = "installation"
	users         = "user"
	follows       = "follow"
	tweets        = "tweet"
)

var mdb = env.GetMongoDB()

// Installations return *Collection for "installation" collection
func Installations() (c *mongo.Collection, err error) { return mdb.GetCollection(installations, true) }

// Users return *Collection for "user" collection
func Users() (c *mongo.Collection, err error) { return mdb.GetCollection(users, true) }

// Follows return *Collection for "follow" collection
func Follows() (c *mongo.Collection, err error) { return mdb.GetCollection(follows, true) }

// Tweets return *Collection for "tweet" collection
func Tweets() (c *mongo.Collection, err error) { return mdb.GetCollection(tweets, true) }
