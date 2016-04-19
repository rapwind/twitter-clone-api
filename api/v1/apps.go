package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/entity"
	"github.com/techcampman/twitter-d-server/errors"
	"github.com/techcampman/twitter-d-server/jsonschema"
	"github.com/techcampman/twitter-d-server/service"
)

func createInstallation(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errors.Send(c, errors.RequestEntityTooLarge())
		return
	}

	// validate
	err = jsonschema.JSONSchema(b, jsonschema.V1CreateInstallationDocument)
	if err != nil {
		errors.Send(c, err)
		return
	}

	// unmarshal json
	i := new(entity.Installation)
	if err = json.Unmarshal(b, i); err != nil {
		errors.Send(c, errors.BadParams("body", string(b)))
		return
	}

	// create installation
	if err = service.CreateInstallation(i); err != nil {
		errors.Send(c, fmt.Errorf("failed to register a installation"))
		return
	}

	// set header
	c.Writer.Header().Set(constant.XPoppoInstallationID, i.UUID)

	c.AbortWithStatus(http.StatusCreated)
}
