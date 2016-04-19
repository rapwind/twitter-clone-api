package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/db/collection"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/jsonschema"
	"github.com/techcampman/twitter-d-server/middleware"
	"github.com/techcampman/twitter-d-server/service"
	"github.com/techcampman/twitter-d-server/utils"
	"gopkg.in/mgo.v2/bson"
)

func signIn(c *gin.Context) {

	// validate request
	installationID := c.Request.Header.Get(constant.XPoppoInstallationID)
	if installationID == "" {
		errors.Send(c, errors.Unauthorized())
		return
	}
	i, err := service.ReadInstallationByUUID(installationID)
	if err != nil {
		errors.Send(c, errors.BadParams(constant.XPoppoInstallationID, installationID))
		return
	}

	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors.Send(c, errors.RequestEntityTooLarge())
		return
	}
	if len(b) == 0 {
		errors.Send(c, errors.Unauthorized())
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

	q := createSignInQueryFromRequest(req)
	u := new(entity.User)
	err = users.Find(q).One(u)
	if err != nil {
		errors.Send(c, errors.Unauthorized())
		return
	}
	ud := new(entity.UserDetail)
	ud.User = u
	if err := service.AppendCounterToUserDetail(ud); err != nil {
		errors.Send(c, err)
	}

	s := new(entity.Session)
	s.UserID = u.ID
	s.InstallationID = i.ID
	middleware.SetSession(c, s)

	c.JSON(http.StatusOK, ud)
}

func signOut(c *gin.Context) {
	// delete auth data from user account
	if err := middleware.DelSession(c); err != nil {
		errors.Send(c, fmt.Errorf("failed to logout a user"))
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func createSignInQueryFromRequest(req *entity.SessionRequest) (q bson.M) {
	emailRegexp := regexp.MustCompile(`[\w.\-]+@[\w\-]+\.[\w.\-]+`)
	if emailRegexp.MatchString(req.AccountName) {
		// if email
		q = bson.M{"email": req.AccountName, "passwordHash": req.PasswordHash}
		return
	}

	phoneRegexp := regexp.MustCompile(`^[\d]+$`)
	numberStr := utils.PhoneNumberNormalization(req.AccountName)
	if phoneRegexp.MatchString(numberStr) {
		// if phoneNumber
		q = bson.M{"phoneNumber": numberStr, "passwordHash": req.PasswordHash}
		return
	}
	// if screenName
	q = bson.M{"screenName": req.AccountName, "passwordHash": req.PasswordHash}
	return
}
