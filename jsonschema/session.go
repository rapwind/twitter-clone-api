package jsonschema

import (
	"github.com/techcampman/twitter-d-server/env"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// V1PostSessionDocument ... json schema for PUT /v1/apps/login
	V1PostSessionDocument *gojsonschema.Schema

	//
	// schema definitions
	//

	v1PostSessionSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"screenName":   screenName,
			"passwordHash": passwordHash,
		},
		"additionalProperties": false,
		"required":             []interface{}{"screenName", "passwordHash"},
	}
)

func initSessionSchema() {

	var err error
	V1PostSessionDocument, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(v1PostSessionSchema))
	env.AssertErrForInit(err)

	return
}
