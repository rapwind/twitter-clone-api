package main

import (
	"runtime"
	"strconv"

	"github.com/dogenzaka/gin-tools/logging"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/api/v1"
	"github.com/techcampman/twitter-d-server/env"
	"github.com/techcampman/twitter-d-server/env/on"
	"github.com/techcampman/twitter-d-server/logger"
	cors "github.com/tommy351/gin-cors"
)

func init() {
	// set go max process
	runtime.GOMAXPROCS(runtime.NumCPU())

	// set gin mode
	switch env.GoEnv() {
	case on.ReleaseEnv:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
}

func main() {

	r := gin.New()

	// Logging, CORS
	r.Use(
		logging.AccessLogger(env.GetAccessLogger()),
		logging.ActivityLogger(env.GetActivityLogger(), nil),
		cors.Middleware(cors.Options{AllowOrigins: []string{"*"}}),
	)

	r.GET("/status", func(c *gin.Context) {})

	if env.GoEnv() != on.ReleaseEnv {
		r.Use(gin.Logger())
	}

	v1.AddV1Endpoints(r)

	// Listen and server
	logger.Info("START POPPO API SERVER: ✔")
	r.Run(":" + strconv.Itoa(env.Port()))
}
