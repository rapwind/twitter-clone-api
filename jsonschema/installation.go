package jsonschema

import (
	"github.com/xeipuuv/gojsonschema"
	"github.com/techcampman/twitter-d-server/env"
)

var (
	// V1CreateInstallationDocument ... json schema for POST /v1/apps
	V1CreateInstallationDocument *gojsonschema.Schema

	//
	// schema definitions
	//

	v1CreateInstallationSchema = map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"clientType": clientType,
			"deviceToken": deviceToken,
		},
		"additionalProperties": false,
		"required":             []interface{}{"clientType"},
	}
)

func initInstallationSchema() {

	var err error
	V1CreateInstallationDocument, err = gojsonschema.NewSchema(gojsonschema.NewGoLoader(v1CreateInstallationSchema))
	env.AssertErrForInit(err)

	return
}
