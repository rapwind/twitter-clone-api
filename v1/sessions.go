package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/jsonschema"
	"github.com/techcampman/twitter-d-server/middleware"
	"gopkg.in/mgo.v2/bson"
)

func signin(c *gin.Context) {

	// validate request
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors.Send(c, errors.RequestEntityTooLarge())
		return
	}
	if len(b) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	if err := jsonschema.JSONSchema(b, jsonschema.V1PostSessionDocument); err != nil {
		errors.Send(c, err)
		return
	}
	req := new(entity.SessionRequest)
	if err = json.Unmarshal(b, req); err != nil {
		errors.Send(c, errors.BadParams("body", fmt.Sprint(b)))
		return
	}

	// get user account
	users, err := collection.Users()
	if err != nil {
		errors.Send(c, err)
		return
	}
	defer users.Close()

	u := new(entity.User)
	err = users.Find(bson.M{"screenName": req.ScreenName, "passwordHash": req.PasswordHash}).One(u)
	if err != nil {
		errors.Send(c, errors.Unauthorized())
		return
	}

	middleware.SetSession(c, u.ID)

	c.JSON(http.StatusOK, u)
}

func signout(c *gin.Context) {
	// delete auth data from user account
	if err := middleware.DelSession(c); err != nil {
		errors.Send(c, fmt.Errorf("failed to logout a user"))
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
