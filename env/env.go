package env

import (
	"io"
	"os"

	"net/url"

	"github.com/techcampman/twitter-d-server/db"
	"github.com/techcampman/twitter-d-server/db/mongo"
	"github.com/techcampman/twitter-d-server/env/on"
	"github.com/techcampman/twitter-d-server/logger"
)

// Environment is an interface for getting objects for the environment.
type Environment interface {
	// Init ... Initialize the environment.
	Init() error

	// GetMongoDB get a session of the MongoDB.  The session should be
	// called `defer session.Close()`.
	GetMongoDB() *mongo.DB

	// GetAccessLogger get access log writer
	GetAccessLogger() io.Writer

	// GetActivityLogger get activity log writer
	GetActivityLogger() io.Writer

	// GoEnv get current environment
	GoEnv() string

	// Port get current port
	Port() int

	// GetCache get a cache store
	GetCache() db.Cache

	// GetStorage get a storage
	GetStorage() db.Storage

	// GetStorageURL get a storage URL
	GetStorageURL() *url.URL
}

// Env initialized poppo-api environment
var env Environment

func init() {

	// get environment
	env = getEnv(os.Getenv("POPPO_ENV"))

	// call to Init() method
	AssertErrForInit(env.Init())

}

func getEnv(s string) Environment {
	switch s {
	case on.ReleaseEnv:
		logger.Info("RELEASE ENVIROMENT MODE: ✔")
		return new(on.Release)
	}
	logger.Info("LOCAL ENVIROMENT MODE: ✔")
	return new(on.Local)
}

// AssertErrForInit if err != nil { panic(err); }
func AssertErrForInit(err error) {
	if err != nil {
		panic(err)
	}
}

// GetMongoDB get a copied session from original
func GetMongoDB() *mongo.DB {
	return env.GetMongoDB()
}

// GetAccessLogger get access log writer
func GetAccessLogger() io.Writer {
	return env.GetAccessLogger()
}

// GetActivityLogger get activity log writer
func GetActivityLogger() io.Writer {
	return env.GetActivityLogger()
}

// GoEnv get current environment
func GoEnv() string {
	return env.GoEnv()
}

// Port get current port
func Port() int {
	return env.Port()
}

// GetCache get a cache storage
func GetCache() db.Cache {
	return env.GetCache()
}

// GetStorage get a storage
func GetStorage() db.Storage {
	return env.GetStorage()
}

// GetStorageURL get a storage URL
func GetStorageURL() *url.URL {
	return env.GetStorageURL()
}
