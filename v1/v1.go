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
		v1.POST("/apps", createInstallation)

		session := v1.Group("/sessions")
		{
			session.POST("/", signin)
			session.DELETE("/", signout)
		}
		users := v1.Group("/users")
		{
			// users.Use(middleware.CheckSession())
			users.GET("/:"+constant.IDKey, getUser)
			users.GET("/:"+constant.IDKey+"/following", getFollowing)
			users.GET("/:"+constant.IDKey+"/follower", getFollower)
		}
		tweets := v1.Group("/tweets")
		{
			tweets.GET("", getTweets)
			tweets.GET("/:"+constant.IDKey, getTweet)
		}
		images := v1.Group("/images")
		{
			images.POST("/", uploadImage)
		}
	}
}
