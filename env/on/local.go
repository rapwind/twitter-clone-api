package on

import (
	"io"
	"os"
	"time"

	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/techcampman/twitter-d-server/bridge"
	"github.com/techcampman/twitter-d-server/db"
	"github.com/techcampman/twitter-d-server/db/mongo"
	"github.com/techcampman/twitter-d-server/db/redis"
	"github.com/techcampman/twitter-d-server/db/s3"
)

// LocalEnv is a signature for POPPO_ENV.
const LocalEnv = "local"

// Local defines configuration for local
type Local struct {
	*mongo.DB
	db.Cache
	db.Storage
	bridge.PushMessage
	accessLogger   io.Writer
	activityLogger io.Writer
	goEnv          string
	port           int
	awsConfig      *aws.Config
	storageURL     *url.URL
}

// Init is a initialize method
func (lo *Local) Init() (err error) {
	// GoEnv and port
	lo.goEnv = LocalEnv
	lo.port = 3000

	// Loggers
	lo.accessLogger = os.Stdout
	lo.activityLogger = os.Stdout

	// MongoDB
	mdb, err := db.NewMongoDB(
		[]string{"localhost:27017"},
		10*time.Second,
		"poppo",
		"",
		"",
		"",
		128,
	)
	if err != nil {
		return err
	}
	lo.DB = mdb

	// Cache - currently Redis
	lo.Cache = redis.NewRedisStore(redis.Config{
		Host:              "localhost:6379",
		Password:          "",
		MaxActive:         10,
		MaxIdle:           5,
		IdleTimeout:       5 * time.Minute,
		DefaultExpiration: 5 * time.Minute,
	})

	// AWS Configurations
	lo.awsConfig = &aws.Config{
		Credentials:            credentials.NewStaticCredentials("AKIAIJNN2OVMBFRGDN6Q", "HOQLaMmo/rEuOunKqqNwtPXg4KWqdOwFfuqH9FMN", ""),
		Endpoint:               aws.String(""),
		Region:                 aws.String("us-east-1"),
		DisableSSL:             aws.Bool(false),
		HTTPClient:             http.DefaultClient,
		LogLevel:               aws.LogLevel(aws.LogOff),
		MaxRetries:             aws.Int(3),
		DisableParamValidation: aws.Bool(false),
	}
	// Storage - currently AWS S3
	bucketName := "static.poppo.me"
	lo.Storage = s3.NewStorageByS3(lo.awsConfig, bucketName)
	lo.storageURL, err = url.Parse("https://static.poppo.me")
	if err != nil {
		return err
	}
	// PushMessage - currently AWS SNS
	lo.PushMessage = bridge.NewPushMessageBySNS(lo.awsConfig)

	return
}

// GetMongoDB get a MongoDB
func (lo *Local) GetMongoDB() *mongo.DB {
	return lo.DB
}

// GetAccessLogger get a access logger
func (lo *Local) GetAccessLogger() io.Writer {
	return lo.accessLogger
}

// GetActivityLogger get a activity logger
func (lo *Local) GetActivityLogger() io.Writer {
	return lo.activityLogger
}

// GoEnv get current environment
func (lo *Local) GoEnv() string {
	return lo.goEnv
}

// Port get current port
func (lo *Local) Port() int {
	return lo.port
}

// GetCache get a cache storage
func (lo *Local) GetCache() db.Cache {
	return lo.Cache
}

// GetStorage get a storage
func (lo *Local) GetStorage() db.Storage {
	return lo.Storage
}

// GetStorageURL get a storage URL
func (lo *Local) GetStorageURL() *url.URL {
	return lo.storageURL
}

// GetPushMessage get a push message service
func (lo *Local) GetPushMessage() bridge.PushMessage {
	return lo.PushMessage
}
