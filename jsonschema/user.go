package jsonschema

import (
	"github.com/techcampman/twitter-d-server/env"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// V1PostUserDocument ... json schema for POST /v1/users
	V1PostUserDocument *gojsonschema.Schema

	//
	// schema definitions
	//

	v1PostUserSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"name":         name,
			"screenName":   screenName,
			"passwordHash": passwordHash,
		},
		"additionalProperties": false,
		"required":             []interface{}{"screenName", "passwordHash"},
	}
)

func initUserSchema() {

	var err error
	V1PostUserDocument, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(v1PostUserSchema))
	env.AssertErrForInit(err)

	return
}
