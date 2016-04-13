package v1

import (
	"github.com/dogenzaka/gin-tools/logging"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/env"
)

// AddV1Endpoints adds entity.
func AddV1Endpoints(r *gin.Engine) {
	v1 := r.Group("/v1",
		logging.AccessLogger(env.GetAccessLogger()),
		logging.ActivityLogger(env.GetActivityLogger(), nil),
	)
	{
		users := v1.Group("/users")
		{
			// users.Use(middleware.CheckSession())
			users.GET("/:"+constant.IDKey, getUser)
		}
		tweets := v1.Group("/tweets")
		{
			tweets.GET("/:"+constant.IDKey, getTweet)
		}
	}
}
