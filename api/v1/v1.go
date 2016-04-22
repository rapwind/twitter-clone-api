package v1

import (
	"github.com/dogenzaka/gin-tools/logging"
	"github.com/dogenzaka/gin-tools/validation"
	"github.com/dogenzaka/gin-tools/validation/validator"
	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/env"
	"github.com/techcampman/twitter-d-server/middleware"
)

// AddV1Endpoints adds entity.
func AddV1Endpoints(r *gin.Engine) {
	v1 := r.Group("/v1",
		logging.AccessLogger(env.GetAccessLogger()),
		logging.ActivityLogger(env.GetActivityLogger(), nil),
	)
	{
		apps := v1.Group("/apps")
		{
			apps.POST("", createInstallation)
			apps.PUT("", updateInstallation)
		}

		sessions := v1.Group("/sessions")
		{
			sessions.POST("", signIn)

			sessions.Use(middleware.CheckSession())
			sessions.DELETE("", signOut)
		}
		users := v1.Group("/users")
		{
			users.POST("", registerUser)

			users.Use(
				middleware.SetLoginUserIDIfNotEmpty(),
				validation.ValidatePathParam(constant.IDKey, validator.ObjectID{}),
				middleware.SetUserOnContext(constant.IDKey),
			)
			users.GET("/:"+constant.IDKey, getUser)
			users.GET("/:"+constant.IDKey+"/tweets", getUserTweets)
			users.GET("/:"+constant.IDKey+"/liked/tweets", getUserLikedTweets)

			users.Use(middleware.CheckSession())
			users.GET("/:"+constant.IDKey+"/following", getFollowing)
			users.GET("/:"+constant.IDKey+"/follower", getFollower)

			users.POST("/:"+constant.IDKey+"/follow", doFollow)
			users.DELETE(":"+constant.IDKey+"/follow", undoFollow)
		}
		tweets := v1.Group("/tweets")
		{
			tweets.Use(middleware.SetLoginUserIDIfNotEmpty())
			tweets.GET("", getTweets)
			tweets.POST("", createTweet)

			tweets.Use(
				validation.ValidatePathParam(constant.IDKey, validator.ObjectID{}),
				middleware.SetTweetOnContext(constant.IDKey),
			)
			tweets.GET("/:"+constant.IDKey, getTweet)

			tweets.Use(middleware.CheckSession())
			tweets.DELETE("/:"+constant.IDKey, deleteTweet)

			tweets.POST("/:"+constant.IDKey+"/like", doLike)
			tweets.DELETE("/:"+constant.IDKey+"/like", undoLike)
		}
		images := v1.Group("/images")
		{
			images.Use(middleware.CheckSession())
			images.POST("", uploadImage)
		}
		notifications := v1.Group("/notifications")
		{
			notifications.Use(middleware.SetLoginUserIDIfNotEmpty())
			notifications.GET("", getNotifications)
		}
	}
}
