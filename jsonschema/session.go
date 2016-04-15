package jsonschema

import (
	"github.com/techcampman/twitter-d-server/env"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// V1PostSessionDocument ... json schema for PSOT /v1/sessions
	V1PostSessionDocument *gojsonschema.Schema

	//
	// schema definitions
	//

	v1PostSessionSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"accountName":  accountName,
			"passwordHash": passwordHash,
		},
		"additionalProperties": false,
		"required":             []interface{}{"passwordHash"},
	}
)

func initSessionSchema() {

	var err error
	V1PostSessionDocument, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(v1PostSessionSchema))
	env.AssertErrForInit(err)

	return
}
