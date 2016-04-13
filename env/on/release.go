package on

import (
	"io"
	"os"
	"time"

	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/techcampman/twitter-d-server/db"
	"github.com/techcampman/twitter-d-server/db/mongo"
	"github.com/techcampman/twitter-d-server/db/redis"
	"github.com/techcampman/twitter-d-server/db/s3"
)

// ReleaseEnv is a signature for POPPO_ENV.
const ReleaseEnv = "release"

// Release defines configuration for Release
type Release struct {
	*mongo.DB
	db.Cache
	db.Storage
	accessLogger   io.Writer
	activityLogger io.Writer
	goEnv          string
	port           int
	awsConfig      *aws.Config
	storageURL     *url.URL
}

// Init is a initialize method
func (rel *Release) Init() (err error) {
	// GoEnv and port
	rel.goEnv = ReleaseEnv
	rel.port = 3000

	// Loggers
	rel.accessLogger = os.Stdout
	rel.activityLogger = os.Stdout

	// MongoDB
	// TODO - Fix Release Config
	mdb, err := db.NewMongoDB(
		[]string{"ds023560.mlab.com:23560"},
		10*time.Second,
		"poppo-mongo-dev",
		"",
		"poppo-administrator",
		"OzamasaIsGod",
		128,
	)
	if err != nil {
		return err
	}
	rel.DB = mdb

	// Cache - currently Redis
	rel.Cache = redis.NewRedisStore(redis.Config{
		Host:              "poppo-cache-set.uktbau.ng.0001.use1.cache.amazonaws.com:6379",
		Password:          "",
		MaxActive:         10,
		MaxIdle:           5,
		IdleTimeout:       5 * time.Minute,
		DefaultExpiration: 5 * time.Minute,
	})

	// AWS Configurations
	rel.awsConfig = &aws.Config{
		Credentials:            credentials.NewStaticCredentials("AKIAJ36RD7B6AM3JWFTA", "CPEmcV/QLmQDnpQvICOoIDC2uhz4Mmq+AE0MsrSv", ""),
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
	rel.Storage = s3.NewStorageByS3(rel.awsConfig, bucketName)
	rel.storageURL, err = url.Parse("https://static.poppo.me")
	if err != nil {
		return err
	}

	return
}

// GetMongoDB get a MongoDB
func (rel *Release) GetMongoDB() *mongo.DB {
	return rel.DB
}

// GetAccessLogger get a access logger
func (rel *Release) GetAccessLogger() io.Writer {
	return rel.accessLogger
}

// GetActivityLogger get a activity logger
func (rel *Release) GetActivityLogger() io.Writer {
	return rel.activityLogger
}

// GoEnv get current environment
func (rel *Release) GoEnv() string {
	return rel.goEnv
}

// Port get current port
func (rel *Release) Port() int {
	return rel.port
}

// GetCache get a cache storage
func (rel *Release) GetCache() db.Cache {
	return rel.Cache
}

// GetStorage get a storage
func (rel *Release) GetStorage() db.Storage {
	return rel.Storage
}

// GetStorageURL get a storage URL
func (rel *Release) GetStorageURL() *url.URL {
	return rel.storageURL
}
